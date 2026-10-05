'use client';

import { useEffect, useRef, useState } from 'react';
import { useDictionary } from '@/components/locale-provider';
import styles from '@/app/game/game.module.css';

type Point = { x: number; y: number };
type Objects = { players: Point[]; monsters: Point[]; npcs: Point[]; stands: Point[] };
type Status = 'connecting' | 'connected' | 'reconnecting' | 'disconnected' | 'unavailable';

function record(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function points(value: unknown): Point[] {
  if (!Array.isArray(value)) return [];
  return value.filter((point): point is Point =>
    record(point) &&
    typeof point.x === 'number' && Number.isFinite(point.x) && point.x >= 0 && point.x < 256 &&
    typeof point.y === 'number' && Number.isFinite(point.y) && point.y >= 0 && point.y < 256,
  );
}

function mapObjects(value: unknown): Objects {
  if (!record(value)) return { players: [], monsters: [], npcs: [], stands: [] };
  return {
    players: points(value.players),
    monsters: points(value.monsters),
    npcs: points(value.npcs),
    stands: points(value.stands),
  };
}

function drawMap(canvas: HTMLCanvasElement, terrain: Point[], objects: Objects) {
  const context = canvas.getContext('2d');
  if (!context) return;
  context.setTransform(1, 0, 0, -1, 0, 512);
  context.fillStyle = '#f8faf7';
  context.fillRect(0, 0, 512, 512);
  const draw = (items: Point[], color: string, size: number) => {
    context.fillStyle = color;
    for (const point of items) context.fillRect(point.x * 2, point.y * 2, size, size);
  };
  draw(terrain, '#244239', 2);
  draw(objects.stands, '#8056a0', 4);
  draw(objects.monsters, '#d04b3c', 4);
  draw(objects.npcs, '#3a9d76', 4);
  draw(objects.players, '#2464c1', 4);
  context.strokeStyle = '#2464c1';
  context.lineWidth = 1;
  for (const point of objects.players) {
    context.beginPath();
    context.arc(point.x * 2, point.y * 2, 10, 0, Math.PI * 2);
    context.stroke();
  }
}

export function GameMapViewer({ maps, websocketUrl }: { maps: string[]; websocketUrl: string | null }) {
  const t = useDictionary();
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const selectedRef = useRef(maps[0]);
  const terrainRef = useRef<Point[]>([]);
  const [selected, setSelected] = useState(maps[0]);
  const [name, setName] = useState(maps[0]);
  const [status, setStatus] = useState<Status>(websocketUrl ? 'connecting' : 'unavailable');
  const [counts, setCounts] = useState({ players: 0, monsters: 0, npcs: 0 });
  const [reconnectKey, setReconnectKey] = useState(0);

  useEffect(() => {
    if (!websocketUrl) return;
    let stopped = false;
    let attempts = 0;
    let timer: ReturnType<typeof setTimeout> | undefined;

    function connect() {
      if (stopped) return;
      setStatus(attempts ? 'reconnecting' : 'connecting');
      const socket = new WebSocket(websocketUrl!);
      socketRef.current = socket;
      socket.onopen = () => {
        attempts = 0;
        setStatus('connected');
        socket.send(JSON.stringify({ action: 'SubscribeMap', in: { name: selectedRef.current } }));
      };
      socket.onmessage = (event) => {
        let message: unknown;
        try { message = JSON.parse(event.data); } catch { return; }
        if (!record(message) || message.action !== 'SubscribeMapReply' || !record(message.out)) return;
        const out = message.out;
        if (out.err) return;
        if (out.name === 'map-name' && typeof out.data === 'string') {
          setName(out.data);
        } else if (out.name === 'map-data') {
          terrainRef.current = points(out.data);
          if (canvasRef.current) drawMap(canvasRef.current, terrainRef.current, mapObjects(null));
        } else if (out.name === 'object') {
          const objects = mapObjects(out.data);
          setCounts({ players: objects.players.length, monsters: objects.monsters.length, npcs: objects.npcs.length });
          if (canvasRef.current) drawMap(canvasRef.current, terrainRef.current, objects);
        }
      };
      socket.onerror = () => socket.close();
      socket.onclose = () => {
        if (stopped) return;
        socketRef.current = null;
        setStatus('disconnected');
        const delay = Math.min(1000 * 2 ** attempts, 10000);
        attempts += 1;
        timer = setTimeout(connect, delay);
      };
    }

    connect();
    return () => {
      stopped = true;
      if (timer) clearTimeout(timer);
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, [websocketUrl, reconnectKey]);

  function selectMap(next: string) {
    selectedRef.current = next;
    setSelected(next);
    setName(next);
    terrainRef.current = [];
    setCounts({ players: 0, monsters: 0, npcs: 0 });
    if (canvasRef.current) drawMap(canvasRef.current, [], mapObjects(null));
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ action: 'SubscribeMap', in: { name: next } }));
    }
  }

  const statusText = {
    connecting: t.mapConnecting,
    connected: t.mapConnected,
    reconnecting: t.mapReconnecting,
    disconnected: t.mapDisconnected,
    unavailable: t.mapConnectionUnavailable,
  }[status];

  return (
    <div className={styles.layout}>
      <section className={styles.viewer} aria-label={t.gameMap}>
        <div className={styles.canvasWrap}>
          <canvas ref={canvasRef} width="512" height="512" aria-label={name} />
        </div>
        <div className={styles.stats}>
          <strong>{name}</strong>
          <span>{t.mapPlayers}: {counts.players}</span>
          <span>{t.mapMonsters}: {counts.monsters}</span>
          <span>{t.mapNpcs}: {counts.npcs}</span>
        </div>
      </section>
      <aside className={styles.controls}>
        <div className={styles.connection}>
          <span className={`${styles.indicator} ${status === 'connected' ? styles.online : ''}`} aria-hidden="true" />
          <span role="status">{statusText}</span>
          {status !== 'connected' && websocketUrl && (
            <button type="button" className={styles.retry} onClick={() => setReconnectKey((key) => key + 1)}>
              {t.reconnect}
            </button>
          )}
        </div>
        <label className={styles.mapLabel} htmlFor="game-map-select">{t.map}</label>
        <select id="game-map-select" className={styles.mapSelect} value={selected} onChange={(event) => selectMap(event.target.value)}>
          {maps.map((map) => <option key={map} value={map}>{map}</option>)}
        </select>
      </aside>
    </div>
  );
}
