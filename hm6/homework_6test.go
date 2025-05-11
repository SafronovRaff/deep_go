package main

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	nameSize   = 42
	finishStr  = 0
	bitsInByte = 8

	shift    = 4
	highMask = 0xF0 // 1111 0000
	lowMask  = 0x0F // 0000 1111

	maskManaHigh   = 0xC0 // 1100 0000
	sixShift       = 6    // сдвиг в право для 0_3
	maskHealthHigh = 0x03 // 0000 0011

	typeMask   = 0xC0 // 1100 0000
	houseFlag  = 0x04 // 0000 0100
	gunFlag    = 0x02 // 0000 0010
	familyFlag = 0x01 // 0000 0001

)

type GamePerson struct {
	name            [nameSize]byte
	strengthRespect byte // [ssssrrrr]
	levelExperience byte // [lllleeee]
	// 42 + 1  + 1
	gold, x, y, z      int32
	healthMana         [3]byte
	typeHouseGunFamily byte // type **_____, House____*__, Gun______*_, family_______*,
	// 16 + 3 + 1
}

func main() {
	fmt.Println("size", unsafe.Sizeof(GamePerson{}))

	p := NewGamePerson(WithName("goog"))
	fmt.Printf("name: %q \n", p.Name())

	p1 := NewGamePerson(WithRespect(9), WithStrength(11))
	fmt.Printf("respect: %d srength: %d\n", p1.Respect(), p1.Strength())

	p = NewGamePerson(WithExperience(8), WithLevel(7))
	fmt.Printf(" exp: %d lvl: %d\n", p.Experience(), p.Level())

	p = NewGamePerson(WithName("goog"), WithCoordinates(17, 12, 22), WithGold(2_000_000))
	fmt.Printf("name: %q xyz %d %d %d gold: %d \n", p.Name(), p.X(), p.Y(), p.Z(), p.Gold())

	p = NewGamePerson(WithMana(900), WithHealth(700))
	fmt.Printf(" mana=%d health=%d \n", p.Mana(), p.Health())

	p2 := NewGamePerson(WithHouse(), WithType(BuilderGamePersonType))
	fmt.Printf("house=%v, gun=%v, family=%v, type=%d \n",
		p2.HasHouse(), p2.HasGun(), p2.HasFamily(), p2.Type())

	p3 := NewGamePerson(WithHouse(), WithGun(), WithFamily(), WithType(WarriorGamePersonType))
	fmt.Printf("house=%v, gun=%v, family=%v, type=%d  \n",
		p3.HasHouse(), p3.HasGun(), p3.HasFamily(), p3.Type())

	p3 = NewGamePerson(WithHouse(), WithGun(), WithFamily(), WithType(BlacksmithGamePersonType))
	fmt.Printf("house=%v, gun=%v, family=%v, type=%d  \n",
		p3.HasHouse(), p3.HasGun(), p3.HasFamily(), p3.Type())

}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, c := range name {
			person.name[i] = byte(c)
		}
		if len(name) < nameSize {
			person.name[len(name)] = finishStr
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthMana[0] &= ^byte(maskManaHigh)
		person.healthMana[0] |= byte((mana >> bitsInByte) << sixShift)
		person.healthMana[1] = byte(mana & 0xFF)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthMana[0] &= ^byte(maskHealthHigh)
		person.healthMana[0] |= byte((health >> bitsInByte) & maskHealthHigh)
		person.healthMana[2] = byte(health & 0xFF)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.strengthRespect &= highMask
		person.strengthRespect |= byte(respect & lowMask)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.strengthRespect &= lowMask
		person.strengthRespect |= byte(strength&lowMask) << shift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExperience &= highMask
		person.levelExperience |= byte(experience & lowMask)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExperience &= lowMask
		person.levelExperience |= byte(level&lowMask) << shift
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= houseFlag
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= gunFlag
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= familyFlag
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily &= ^byte(typeMask)
		person.typeHouseGunFamily |= byte(personType) << sixShift
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

func NewGamePerson(options ...Option) GamePerson {
	res := GamePerson{}
	for _, option := range options {
		option(&res)
	}

	return res
}

func (p *GamePerson) Name() string {
	i := 0
	for ; i < nameSize; i++ {
		if p.name[i] == finishStr {
			break
		}
	}
	return string(p.name[:i])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	low := int(p.healthMana[1])
	top := int((p.healthMana[0] & maskManaHigh) >> sixShift)
	return (top << bitsInByte) | low
}

func (p *GamePerson) Health() int {
	low := int(p.healthMana[2])
	top := int(p.healthMana[0] & 0x03)
	return (top << bitsInByte) | low
}

func (p *GamePerson) Respect() int {
	return int(p.strengthRespect & lowMask)
}

func (p *GamePerson) Strength() int {
	return int(p.strengthRespect>>shift) & lowMask
}

func (p *GamePerson) Experience() int {
	return int(p.levelExperience & lowMask)
}

func (p *GamePerson) Level() int {
	return int(p.levelExperience >> shift)
}

func (p *GamePerson) HasHouse() bool {
	return (p.typeHouseGunFamily & houseFlag) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.typeHouseGunFamily & gunFlag) != 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.typeHouseGunFamily & familyFlag) != 0
}

func (p *GamePerson) Type() int {
	return int(p.typeHouseGunFamily >> sixShift)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
