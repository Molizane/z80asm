package z80asm

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
)

var (
	basicPlaneTable = []string{
		"nop", "ld bc, **", "ld (bc), a", "inc bc", "inc b", "dec b", "ld b, *", "rlca",
		"ex af, af'", "add hl, bc", "ld a, (bc)", "dec bc", "inc c", "dec c", "ld c, *", "rrca",
		"djnz *", "ld de, **", "ld (de), a", "inc de", "inc d", "dec d", "ld d, *", "rla",
		"jr *", "add hl, de", "ld a, (de)", "dec de", "inc e", "dec e", "ld e, *", "rra",
		"jr nz, *", "ld hl, **", "ld (**), hl", "inc hl", "inc h", "dec h", "ld h, *", "daa",
		"jr z, *", "add hl, hl", "ld hl, (**)", "dec hl", "inc l", "dec l", "ld l, *", "cpl",
		"jr nc, *", "ld sp, **", "ld (**), a", "inc sp", "inc (hl)", "dec (hl)", "ld (hl), *", "scf",
		"jr c, *", "add hl, sp", "ld a, (**)", "dec sp", "inc a", "dec a", "ld a, *", "ccf",
		"ld b, b", "ld b, c", "ld b, d", "ld b, e", "ld b, h", "ld b, l", "ld b, (hl)", "ld b, a",
		"ld c, b", "ld c, c", "ld c, d", "ld c, e", "ld c, h", "ld c, l", "ld c, (hl)", "ld c, a",
		"ld d, b", "ld d, c", "ld d, d", "ld d, e", "ld d, h", "ld d, l", "ld d, (hl)", "ld d, a",
		"ld e, b", "ld e, c", "ld e, d", "ld e, e", "ld e, h", "ld e, l", "ld e, (hl)", "ld e, a",
		"ld h, b", "ld h, c", "ld h, d", "ld h, e", "ld h, h", "ld h, l", "ld h, (hl)", "ld h, a",
		"ld l, b", "ld l, c", "ld l, d", "ld l, e", "ld l, h", "ld l, l", "ld l, (hl)", "ld l, a",
		"ld (hl), b", "ld (hl), c", "ld (hl), d", "ld (hl), e", "ld (hl), h", "ld (hl), l", "halt", "ld (hl), a",
		"ld a, b", "ld a, c", "ld a, d", "ld a, e", "ld a, h", "ld a, l", "ld a, (hl)", "ld a, a",
		"add a, b", "add a, c", "add a, d", "add a, e", "add a, h", "add a, l", "add a, (hl)", "add a, a",
		"adc a, b", "adc a, c", "adc a, d", "adc a, e", "adc a, h", "adc a, l", "adc a, (hl)", "adc a, a",
		"sub b", "sub c", "sub d", "sub e", "sub h", "sub l", "sub (hl)", "sub a",
		"sbc a, b", "sbc a, c", "sbc a, d", "sbc a, e", "sbc a, h", "sbc a, l", "sbc a, (hl)", "sbc a, a",
		"and b", "and c", "and d", "and e", "and h", "and l", "and (hl)", "and a",
		"xor b", "xor c", "xor d", "xor e", "xor h", "xor l", "xor (hl)", "xor a",
		"or b", "or c", "or d", "or e", "or h", "or l", "or (hl)", "or a",
		"cp b", "cp c", "cp d", "cp e", "cp h", "cp l", "cp (hl)", "cp a",
		"ret nz", "pop bc", "jp nz, **", "jp **", "call nz, **", "push bc", "add a, *", "rst 0",
		"ret z", "ret", "jp z, **", "", "call z, **", "call **", "adc a, *", "rst 0x08",
		"ret nc", "pop de", "jp nc, **", "out (*), a", "call nc, **", "push de", "sub *", "rst 0x10",
		"ret c", "exx", "jp c, **", "in a, (*)", "call c, **", "", "sbc a, *", "rst 0x18",
		"ret po", "pop hl", "jp po, **", "ex (sp), hl", "call po, **", "push hl", "and *", "rst 0x20",
		"ret pe", "jp (hl)", "jp pe, **", "ex de, hl", "call pe, **", "", "xor *", "rst 0x28",
		"ret p", "pop af", "jp p, **", "di", "call p, **", "push af", "or *", "rst 0x30",
		"ret m", "ld sp, hl", "jp m, **", "ei", "call m, **", "", "cp *", "rst 0x38",
	}

	// non-standard, undocumented, or the second and subsequent dupes are
	// prefixed with "?"
	extendedPlaneTable = []string{
		"in b, (c)", "out (c), b", "sbc hl, bc", "ld (**), bc", "neg", "retn", "im 0", "ld i, a",
		"in c, (c)", "out (c), c", "adc hl, bc", "ld bc, (**)", "?neg", "reti", "?im 0/1", "ld r, a",
		"in d, (c)", "out (c), d", "sbc hl, de", "ld (**), de", "?neg", "?retn", "im 1", "ld a, i",
		"in e, (c)", "out (c), e", "adc hl, de", "ld de, (**)", "?neg", "?retn", "im 2", "ld a, r",
		"in h, (c)", "out (c), h", "sbc hl, hl", "?ld (**), hl", "?neg", "?retn", "?im 0", "rrd",
		"in l, (c)", "out (c), l", "adc hl, hl", "?ld hl, (**)", "?neg", "?retn", "?im 0/1", "rld",
		"?in (c)", "?out (c), 0", "sbc hl, sp", "ld (**), sp", "?neg", "?retn", "?im 1", "",
		"in a, (c)", "out (c), a", "adc hl, sp", "ld sp, (**)", "?neg", "?retn", "?im 2", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"ldi", "cpi", "ini", "outi", "", "", "", "",
		"ldd", "cpd", "ind", "outd", "", "", "", "",
		"ldir", "cpir", "inir", "otir", "", "", "", "",
		"lddr", "cpdr", "indr", "otdr", "", "", "", "",
	}

	bitPlaneTable = []string{
		"rlc b", "rlc c", "rlc d", "rlc e", "rlc h", "rlc l", "rlc (hl)", "rlc a",
		"rrc b", "rrc c", "rrc d", "rrc e", "rrc h", "rrc l", "rrc (hl)", "rrc a",
		"rl b", "rl c", "rl d", "rl e", "rl h", "rl l", "rl (hl)", "rl a",
		"rr b", "rr c", "rr d", "rr e", "rr h", "rr l", "rr (hl)", "rr a",
		"sla b", "sla c", "sla d", "sla e", "sla h", "sla l", "sla (hl)", "sla a",
		"sra b", "sra c", "sra d", "sra e", "sra h", "sra l", "sra (hl)", "sra a",
		"?sll b", "?sll c", "?sll d", "?sll e", "?sll h", "?sll l", "?sll (hl)", "?sll a",
		"srl b", "srl c", "srl d", "srl e", "srl h", "srl l", "srl (hl)", "srl a",
		"bit 0, b", "bit 0, c", "bit 0, d", "bit 0, e", "bit 0, h", "bit 0, l", "bit 0, (hl)", "bit 0, a",
		"bit 1, b", "bit 1, c", "bit 1, d", "bit 1, e", "bit 1, h", "bit 1, l", "bit 1, (hl)", "bit 1, a",
		"bit 2, b", "bit 2, c", "bit 2, d", "bit 2, e", "bit 2, h", "bit 2, l", "bit 2, (hl)", "bit 2, a",
		"bit 3, b", "bit 3, c", "bit 3, d", "bit 3, e", "bit 3, h", "bit 3, l", "bit 3, (hl)", "bit 3, a",
		"bit 4, b", "bit 4, c", "bit 4, d", "bit 4, e", "bit 4, h", "bit 4, l", "bit 4, (hl)", "bit 4, a",
		"bit 5, b", "bit 5, c", "bit 5, d", "bit 5, e", "bit 5, h", "bit 5, l", "bit 5, (hl)", "bit 5, a",
		"bit 6, b", "bit 6, c", "bit 6, d", "bit 6, e", "bit 6, h", "bit 6, l", "bit 6, (hl)", "bit 6, a",
		"bit 7, b", "bit 7, c", "bit 7, d", "bit 7, e", "bit 7, h", "bit 7, l", "bit 7, (hl)", "bit 7, a",
		"res 0, b", "res 0, c", "res 0, d", "res 0, e", "res 0, h", "res 0, l", "res 0, (hl)", "res 0, a",
		"res 1, b", "res 1, c", "res 1, d", "res 1, e", "res 1, h", "res 1, l", "res 1, (hl)", "res 1, a",
		"res 2, b", "res 2, c", "res 2, d", "res 2, e", "res 2, h", "res 2, l", "res 2, (hl)", "res 2, a",
		"res 3, b", "res 3, c", "res 3, d", "res 3, e", "res 3, h", "res 3, l", "res 3, (hl)", "res 3, a",
		"res 4, b", "res 4, c", "res 4, d", "res 4, e", "res 4, h", "res 4, l", "res 4, (hl)", "res 4, a",
		"res 5, b", "res 5, c", "res 5, d", "res 5, e", "res 5, h", "res 5, l", "res 5, (hl)", "res 5, a",
		"res 6, b", "res 6, c", "res 6, d", "res 6, e", "res 6, h", "res 6, l", "res 6, (hl)", "res 6, a",
		"res 7, b", "res 7, c", "res 7, d", "res 7, e", "res 7, h", "res 7, l", "res 7, (hl)", "res 7, a",
		"set 0, b", "set 0, c", "set 0, d", "set 0, e", "set 0, h", "set 0, l", "set 0, (hl)", "set 0, a",
		"set 1, b", "set 1, c", "set 1, d", "set 1, e", "set 1, h", "set 1, l", "set 1, (hl)", "set 1, a",
		"set 2, b", "set 2, c", "set 2, d", "set 2, e", "set 2, h", "set 2, l", "set 2, (hl)", "set 2, a",
		"set 3, b", "set 3, c", "set 3, d", "set 3, e", "set 3, h", "set 3, l", "set 3, (hl)", "set 3, a",
		"set 4, b", "set 4, c", "set 4, d", "set 4, e", "set 4, h", "set 4, l", "set 4, (hl)", "set 4, a",
		"set 5, b", "set 5, c", "set 5, d", "set 5, e", "set 5, h", "set 5, l", "set 5, (hl)", "set 5, a",
		"set 6, b", "set 6, c", "set 6, d", "set 6, e", "set 6, h", "set 6, l", "set 6, (hl)", "set 6, a",
		"set 7, b", "set 7, c", "set 7, d", "set 7, e", "set 7, h", "set 7, l", "set 7, (hl)", "set 7, a",
	}

	ixPlaneTable = []string{
		"", "", "", "", "", "", "", "", "", "add ix, bc", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "add ix, de", "", "", "", "", "", "",
		"", "ld ix, **", "ld (**), ix", "inc ix", "?inc ixh", "?dec ixh", "?ld ixh, *", "", "", "add ix, ix", "ld ix, (**)", "dec ix", "?inc ixl", "?dec ixl", "?ld ixl, *", "",
		"", "", "", "", "inc (ix+*)", "dec (ix+*)", "ld (ix+*), *", "", "", "add ix, sp", "", "", "", "", "", "",
		"", "", "", "", "?ld b, ixh", "?ld b, ixl", "ld b, (ix+*)", "", "", "", "", "", "?ld c, ixh", "?ld c, ixl", "ld c, (ix+*)", "",
		"", "", "", "", "?ld d, ixh", "?ld d, ixl", "ld d, (ix+*)", "", "", "", "", "", "?ld e, ixh", "?ld e, ixl", "ld e, (ix+*)", "",
		"?ld ixh, b", "?ld ixh, c", "?ld ixh, d", "?ld ixh, e", "?ld ixh, ixh", "?ld ixh, ixl", "ld h, (ix+*)", "?ld ixh, a", "?ld ixl, b", "?ld ixl, c", "?ld ixl, d", "?ld ixl, e", "?ld ixl, ixh", "?ld ixl, ixl", "ld l, (ix+*)", "?ld ixl, a",
		"ld (ix+*), b", "ld (ix+*), c", "ld (ix+*), d", "ld (ix+*), e", "ld (ix+*), h", "ld (ix+*), l", "", "ld (ix+*), a", "", "", "", "", "?ld a, ixh", "?ld a, ixl", "ld a, (ix+*)", "",
		"", "", "", "", "?add a, ixh", "?add a, ixl", "add a, (ix+*)", "", "", "", "", "", "?adc a, ixh", "?adc a, ixl", "adc a, (ix+*)", "",
		"", "", "", "", "?sub ixh", "?sub ixl", "sub (ix+*)", "", "", "", "", "", "?sbc a, ixh", "?sbc a, ixl", "sbc a, (ix+*)", "",
		"", "", "", "", "?and ixh", "?and ixl", "and (ix+*)", "", "", "", "", "", "?xor ixh", "?xor ixl", "xor (ix+*)", "",
		"", "", "", "", "?or ixh", "?or ixl", "or (ix+*)", "", "", "", "", "", "?cp ixh", "?cp ixl", "cp (ix+*)", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "pop ix", "", "ex (sp), ix", "", "push ix", "", "", "", "jp (ix)", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "ld sp, ix", "", "", "", "", "", "",
	}
	iyPlaneTable = []string{
		"", "", "", "", "", "", "", "", "", "add iy, bc", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "add iy, de", "", "", "", "", "", "",
		"", "ld iy, **", "ld (**), iy", "inc iy", "?inc iyh", "?dec iyh", "?ld iyh, *", "", "", "add iy, iy", "ld iy, (**)", "dec iy", "?inc iyl", "?dec iyl", "?ld iyl, *", "",
		"", "", "", "", "inc (iy+*)", "dec (iy+*)", "ld (iy+*), *", "", "", "add iy, sp", "", "", "", "", "", "",
		"", "", "", "", "?ld b, iyh", "?ld b, iyl", "ld b, (iy+*)", "", "", "", "", "", "?ld c, iyh", "?ld c, iyl", "ld c, (iy+*)", "",
		"", "", "", "", "?ld d, iyh", "?ld d, iyl", "ld d, (iy+*)", "", "", "", "", "", "?ld e, iyh", "?ld e, iyl", "ld e, (iy+*)", "",
		"?ld iyh, b", "?ld iyh, c", "?ld iyh, d", "?ld iyh, e", "?ld iyh, iyh", "?ld iyh, iyl", "ld h, (iy+*)", "?ld iyh, a", "?ld iyl, b", "?ld iyl, c", "?ld iyl, d", "?ld iyl, e", "?ld iyl, iyh", "?ld iyl, iyl", "ld l, (iy+*)", "?ld iyl, a",
		"ld (iy+*), b", "ld (iy+*), c", "ld (iy+*), d", "ld (iy+*), e", "ld (iy+*), h", "ld (iy+*), l", "", "ld (iy+*), a", "", "", "", "", "?ld a, iyh", "?ld a, iyl", "ld a, (iy+*)", "",
		"", "", "", "", "?add a, iyh", "?add a, iyl", "add a, (iy+*)", "", "", "", "", "", "?adc a, iyh", "?adc a, iyl", "adc a, (iy+*)", "",
		"", "", "", "", "?sub iyh", "?sub iyl", "sub (iy+*)", "", "", "", "", "", "?sbc a, iyh", "?sbc a, iyl", "sbc a, (iy+*)", "",
		"", "", "", "", "?and iyh", "?and iyl", "and (iy+*)", "", "", "", "", "", "?xor iyh", "?xor iyl", "xor (iy+*)", "",
		"", "", "", "", "?or iyh", "?or iyl", "or (iy+*)", "", "", "", "", "", "?cp iyh", "?cp iyl", "cp (iy+*)", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "pop iy", "", "ex (sp), iy", "", "push iy", "", "", "", "jp (iy)", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "ld sp, iy", "", "", "", "", "", "",
	}
	ixBitPlaneTable = padOut([]string{
		"rlc (ix+*)", "rrc (ix+*)",
		"rl (ix+*)", "rr (ix+*)",
		"sla (ix+*)", "sra (ix+*)",
		"", "srl (ix+*)",
		"bit 0, (ix+*)", "bit 1, (ix+*)",
		"bit 2, (ix+*)", "bit 3, (ix+*)",
		"bit 4, (ix+*)", "bit 5, (ix+*)",
		"bit 6, (ix+*)", "bit 7, (ix+*)",
		"res 0, (ix+*)", "res 1, (ix+*)",
		"res 2, (ix+*)", "res 3, (ix+*)",
		"res 4, (ix+*)", "res 5, (ix+*)",
		"res 6, (ix+*)", "res 7, (ix+*)",
		"set 0, (ix+*)", "set 1, (ix+*)",
		"set 2, (ix+*)", "set 3, (ix+*)",
		"set 4, (ix+*)", "set 5, (ix+*)",
		"set 6, (ix+*)", "set 7, (ix+*)",
	})
	iyBitPlaneTable = padOut([]string{
		"rlc (iy+*)", "rrc (iy+*)",
		"rl (iy+*)", "rr (iy+*)",
		"sla (iy+*)", "sra (iy+*)",
		"", "srl (iy+*)",
		"bit 0, (iy+*)", "bit 1, (iy+*)",
		"bit 2, (iy+*)", "bit 3, (iy+*)",
		"bit 4, (iy+*)", "bit 5, (iy+*)",
		"bit 6, (iy+*)", "bit 7, (iy+*)",
		"res 0, (iy+*)", "res 1, (iy+*)",
		"res 2, (iy+*)", "res 3, (iy+*)",
		"res 4, (iy+*)", "res 5, (iy+*)",
		"res 6, (iy+*)", "res 7, (iy+*)",
		"set 0, (iy+*)", "set 1, (iy+*)",
		"set 2, (iy+*)", "set 3, (iy+*)",
		"set 4, (iy+*)", "set 5, (iy+*)",
		"set 6, (iy+*)", "set 7, (iy+*)",
	})
)

var planeTestTables = []struct {
	prefix     []byte
	table      []string
	start, end int
}{
	{
		table: basicPlaneTable,
	},
	{
		prefix: []byte{0xed},
		table:  extendedPlaneTable,
		start:  0x40,
		end:    0xc0,
	},
	{
		prefix: []byte{0xcb},
		table:  bitPlaneTable,
	},
	{
		prefix: []byte{0xdd},
		table:  ixPlaneTable,
	},
	{
		prefix: []byte{0xfd},
		table:  iyPlaneTable,
	},
	{
		prefix: []byte{0xdd, 0xcb},
		table:  ixBitPlaneTable,
	},
	{
		prefix: []byte{0xfd, 0xcb},
		table:  iyBitPlaneTable,
	},
}

func padOut(xs []string) []string {
	var r []string
	k := 0
	for i := 0; i < 256; i++ {
		if i%16 == 6 || i%16 == 14 {
			r = append(r, xs[k])
			k++
		} else {
			r = append(r, "")
		}
	}
	return r
}

func planeDiff(got, want []string, start, end int) error {
	if end == 0 {
		end = 256
	}
	if len(got) != 256 {
		return fmt.Errorf("got %d instructions, want %d", len(got), 256)
	}
	if len(want) != end-start {
		return fmt.Errorf("error in test spec: %d instructions, want %d", len(want), end-start)
	}
	var errs []string
	for i := 0; i < 256; i++ {
		if i < start || i >= end {
			if got[i] != "" {
				errs = append(errs, fmt.Sprintf("%02x: got %q, expected nothing", i, got[i]))
			}
			continue
		}
		w := want[i-start]
		if len(w) > 0 && w[0] == '?' {
			w = ""
		}
		if w != got[i] {
			errs = append(errs, fmt.Sprintf("%02x: got %q, want %q", i, got[i], w))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// getPlane returns a map of instructions with the given prefix.
func getPlane(asm *Assembler, prefix []byte) []string {
	result := make([]string, 256)
	collisions := make([][]string, 256)
	for cmd, asm := range asm.commandTable {
		switch v := asm.(type) {
		case commandAssembler:
			for o, bs := range v.args {
				if b, ok := getByte(prefix, bs); ok {
					s := fmt.Sprintf("%s %s", cmd, o)
					if o == void {
						s = cmd
					}
					result[b] = s
					collisions[b] = append(collisions[b], result[b])
				}
			}
		}
	}
	failed := false
	for i, c := range collisions {
		if len(c) > 1 {
			fmt.Printf("collisions at 0x%02x: %s\n", i, strings.Join(c, "; "))
			failed = true
		}
	}
	if failed {
		log.Fatalf("found collisions!")
	}
	return result
}

func TestPlanes(t *testing.T) {
	asm, err := NewAssembler(nil)
	if err != nil {
		t.Fatalf("failed to create assembler: %v", err)
	}
	for _, tc := range planeTestTables {
		p := getPlane(asm, tc.prefix)
		if err := planeDiff(p, tc.table, tc.start, tc.end); err != nil {
			t.Errorf("Instructions for prefix %v differ:\n%v", tc.prefix, err)
		}
	}
}
