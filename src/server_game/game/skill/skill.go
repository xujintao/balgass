package skill

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Skill struct {
	*SkillBase       `json:"-"`
	*MasterSkillBase `json:"-"`
	Index            int `json:"index"`
	MasterClass      int `json:"master_class,omitempty"`
	MasterType       int `json:"master_type,omitempty"`
	MasterIndex      int `json:"master_index,omitempty"`
	Level            int `json:"level,omitempty"`
	DamageMin        int `json:"-"`
	DamageMax        int `json:"-"`
}

type Skills map[int]*Skill

func (s Skills) MarshalJSON() ([]byte, error) {
	var skills []*Skill
	for _, v := range s {
		skills = append(skills, v)
	}
	// sort
	data, err := json.Marshal(skills)
	if err != nil {
		return nil, err
	}
	return data, err
}

// func (s Skills) UnmarshalJSON(buf []byte) error {
// 	var skills []*Skill
// 	err := json.Unmarshal(buf, &skills)
// 	if err != nil {
// 		return err
// 	}
// 	s = make(Skills)
// 	for _, v := range skills {
// 		switch {
// 		case v.Index < 300:
// 			base, ok := SkillManager.skillTable[v.Index]
// 			if !ok {
// 				log.Printf("Skills UnmarshalJSON failed skillTable [index]%d\n", v.Index)
// 				continue
// 			}
// 			v.SkillBase = base
// 		default:
// 			index := v.MasterIndex%36 - 1
// 			rank, pos := index>>2, index%4
// 			base := SkillManager.masterSkillTable[v.MasterClass][v.MasterType][rank][pos]
// 			if base == nil {
// 				log.Printf("Skills UnmarshalJSON failed masterSkillTable [class]%d [type]%d [index]%d\n",
// 					v.MasterClass, v.MasterType, v.MasterIndex)
// 				continue
// 			}
// 			v.MasterSkillBase = base
// 			v.Index = base.SkillID
// 		}
// 		s[v.Index] = v
// 	}
// 	return nil
// }

func (s Skills) Value() (driver.Value, error) {
	return s.MarshalJSON()
}

func (s *Skills) Scan(value any) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to Scan Inventory value:", value))
	}
	var skills []*Skill
	err := json.Unmarshal(buf, &skills)
	if err != nil {
		return err
	}
	ts := make(map[int]*Skill)
	for _, v := range skills {
		switch {
		case v.Index < 300:
			base, ok := SkillManager.skillTable[v.Index]
			if !ok {
				log.Printf("Skills UnmarshalJSON failed skillTable [index]%d\n", v.Index)
				continue
			}
			v.SkillBase = base
		default:
			index := v.MasterIndex%36 - 1
			rank, pos := index>>2, index%4
			base := SkillManager.masterSkillTable[v.MasterClass][v.MasterType][rank][pos]
			if base == nil {
				log.Printf("Skills UnmarshalJSON failed masterSkillTable [class]%d [type]%d [index]%d\n",
					v.MasterClass, v.MasterType, v.MasterIndex)
				continue
			}
			v.MasterSkillBase = base
			v.Index = base.SkillID
		}
		ts[v.Index] = v
	}
	*s = ts
	return nil
}

var poolSkill = sync.Pool{
	New: func() interface{} {
		return &Skill{}
	},
}

// Get get a skill from pool
func (s Skills) Get(index, level int) *Skill {
	damage := 0
	if index >= 300 {
		damage = s.getMasterSkillDamage(index, level)
	} else {
		skillBase := SkillManager.skillTable[index]
		damage = skillBase.Damage
	}
	ss := poolSkill.Get().(*Skill)
	ss.Level = level
	ss.DamageMin = damage
	ss.DamageMax = damage + damage/2
	return ss
}

// Put put a skill to pool
func (s Skills) Put(skill *Skill) {
	poolSkill.Put(skill)
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
