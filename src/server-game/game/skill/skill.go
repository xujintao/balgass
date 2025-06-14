package skill

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
)

func init() {
	Skill0 = &Skill{
		SkillBase: &SkillBase{},
	}
}

var Skill0 *Skill

type Skill struct {
	*SkillBase `json:"-"`
	Index      int     `json:"index"`
	Level      int     `json:"level,omitempty"`
	DamageMin  int     `json:"-"`
	DamageMax  int     `json:"-"`
	UIIndex    int     `json:"-"`
	CurValue   float32 `json:"-"`
	NextValue  float32 `json:"-"`
}

type SortedSkillSlice []*Skill

func (s SortedSkillSlice) Len() int {
	return len(s)
}

func (s SortedSkillSlice) Less(i, j int) bool {
	return s[i].Index < s[j].Index
}

func (s SortedSkillSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *Skill) Marshal() ([]byte, error) {
	var bw bytes.Buffer
	binary.Write(&bw, binary.LittleEndian, uint16(s.Index))
	level := s.Level << 3
	level |= int(s.Index) & 0x07
	bw.WriteByte(byte(level))
	return bw.Bytes(), nil
}

type Skills map[int]*Skill

func (s Skills) MarshalJSON() ([]byte, error) {
	var skills []*Skill
	for _, v := range s {
		skills = append(skills, v)
	}
	sort.Sort(SortedSkillSlice(skills))
	data, err := json.Marshal(skills)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *Skills) UnmarshalJSON(buf []byte) error {
	var skills []*Skill
	err := json.Unmarshal(buf, &skills)
	if err != nil {
		return err
	}
	ts := make(Skills)
	for _, v := range skills {
		skillBase, ok := SkillManager.skillTable[v.Index]
		if !ok {
			log.Printf("Skills UnmarshalJSON failed skillTable [index]%d\n", v.Index)
			continue
		}
		v.SkillBase = skillBase
		ts[v.Index] = v
	}
	*s = ts
	return nil
}

func (s Skills) Value() (driver.Value, error) {
	return s.MarshalJSON()
}

func (s *Skills) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Skills value:", value))
	}
	return s.UnmarshalJSON(buf)
}

// Get get a skill from pool
func (s Skills) Get(index int) (*Skill, bool) {
	skillBase, ok := SkillManager.skillTable[index]
	if !ok {
		return nil, false
	}
	_, ok = s[index]
	if ok {
		return nil, false
	}
	damage := skillBase.Damage
	ss := &Skill{
		SkillBase: skillBase,
		Index:     index,
		Level:     0,
		DamageMin: damage,
		DamageMax: damage + damage/2,
	}
	s[index] = ss
	return ss, true
}

func (s Skills) GetMaster(class, index, point int, f func(point, uiIndex, index, level int, curValue, NextValue float32)) bool {
	skillBase, ok := SkillManager.skillTable[index]
	if !ok {
		return false
	}
	masterSkillBase, ok := SkillManager.getMasterSkillBase(class, index)
	if !ok {
		return false
	}
	pSkill1 := masterSkillBase.ParentSkill1
	if pSkill1 != 0 {
		_, ok := SkillManager.skillTable[pSkill1]
		if !ok {
			return false
		}
	}
	pSkill2 := masterSkillBase.ParentSkill2
	if pSkill2 != 0 {
		_, ok := SkillManager.skillTable[pSkill2]
		if !ok {
			return false
		}
	}
	if point < masterSkillBase.ReqMinPoint {
		return false
	}
	ss, ok := s[index]
	if ok {
		if ss.Level+masterSkillBase.ReqMinPoint > 20 {
			return false
		}
		ss.Level += masterSkillBase.ReqMinPoint
	} else {
		ss = &Skill{
			SkillBase: skillBase,
			Index:     index,
			UIIndex:   masterSkillBase.Index,
			Level:     masterSkillBase.ReqMinPoint,
			CurValue:  0,
			NextValue: 0,
		}
		s[index] = ss
	}
	// if index >= 300 {
	// 	damage = s.getMasterSkillDamage(index, level)
	// }
	f(masterSkillBase.ReqMinPoint, masterSkillBase.Index, int(index), ss.Level, ss.CurValue, ss.NextValue)
	return true
}

func (s Skills) Put(index int) (*Skill, bool) {
	ss, ok := s[index]
	if !ok {
		return nil, false
	}
	delete(s, index)
	return ss, true
}

// example 1:
// skill 342 use 336+43
// 342->339
// 339->336
// 336->43

// example 2:
// skill 339 use 336+43
// 339->336
// 336->43

// example 3:
// skill 336 use 336+43
// 336->43

// example 4:
// skill 346 use 346
// 346->344
// 344->0
func (s Skills) getMasterSkillDamage(index, level int) int {
	brand1 := index
	brand2 := index
	for i := 0; i < 3; i++ {
		brand1 = brand2
		brand2 = SkillManager.skillTable[brand2].Brand
		if brand2 < 300 {
			break
		}
	}
	damage := 0
	if brand1 != index {
		index = brand1
		level = s[index].Level
	}
	stid := SkillManager.skillTable[index].STID
	if SkillManager.masterSkillValueTable[stid].valueType == valueTypeDamage {
		damage += int(SkillManager.masterSkillValueTable[stid].values[level])
	}
	if brand2 > 0 {
		damage += SkillManager.skillTable[brand2].Damage
	}
	return damage
}

func (s Skills) FillSkillData(class int) {
	for _, ss := range s {
		skillBase, ok := SkillManager.skillTable[ss.Index]
		if !ok {
			continue
		}
		switch {
		case ss.Index < 300:
			damage := skillBase.Damage
			ss.DamageMin = damage
			ss.DamageMax = damage + damage/2
		default:
			masterSkillBase, ok := SkillManager.getMasterSkillBase(class, ss.Index)
			if !ok {
				continue
			}
			ss.UIIndex = masterSkillBase.Index
			ss.CurValue = 0
			ss.NextValue = 0
		}
	}
}

func (s Skills) ForEachActiveSkill(f func(*Skill)) {
	for _, ss := range s {
		if ss.UseType == 3 {
			continue
		}
		f(ss)
	}
}

func (s Skills) ForEachMasterSkill(f func(index, level int, curValue, nextValue float32)) {
	for _, ss := range s {
		f(ss.UIIndex, ss.Level, ss.CurValue, ss.NextValue)
	}
}
