package hasher

func rightrot(a string, k int) string {
	return a[len(a)-k:] + a[:len(a)-k]
}

func rightshift(a string, k int) string {
	for i := 0; i < k; i++ {
		a = "0" + a
	}
	return a[:len(a)-k]
}

func add32(a ...string) string {
	ln := 32

	for i := 0; i < len(a); i++ {
		for len(a[i]) < ln {
			a[i] = "0" + a[i]
		}
	}

	res := ""

	ext := 0
	for j := ln - 1; j >= 0; j-- {
		bl := ext

		for _, i := range a {
			if i[j] == '1' {
				bl++
			}
		}

		res = [2]string{"0", "1"}[bl&1] + res
		ext = (bl >> 1)
	}

	return res
}

func and32(a ...string) string {
	ln := 32

	for i := 0; i < len(a); i++ {
		for len(a[i]) < ln {
			a[i] = "0" + a[i]
		}
	}

	res := ""
	for j := 0; j < 32; j++ {
		ok := 1
		for _, i := range a {
			if i[j] == '0' {
				ok = 0
			}
		}

		res += [2]string{"0", "1"}[ok]
	}

	return res
}

func xor32(a ...string) string {
	ln := 32

	for _, el := range a {
		ln = max(ln, len(el))
	}

	for i := 0; i < len(a); i++ {
		for len(a[i]) < ln {
			a[i] = "0" + a[i]
		}
	}

	res := ""

	for j := 0; j < ln; j++ {
		bl := 0
		for _, i := range a {
			if i[j] == '1' {
				bl ^= 1
			}
		}

		res += [2]string{"0", "1"}[bl]
	}

	return res
}

func bin(x uint32) string {
	res := ""

	for i := 31; i >= 0; i-- {
		res += [2]string{"0", "1"}[(x>>i)&1]
	}

	return res
}

func CalcSha256(a string) string {
	r := ""

	for _, e := range a {
		x := uint8(e)

		for j := 7; j >= 0; j-- {
			r += []string{"0", "1"}[(x>>j)&1]
		}
	}

	for len(r)+16 < 512 {
		r += "0"
	}

	for j := 15; j >= 0; j-- {
		r += []string{"0", "1"}[(len(a)>>j)&1]
	}

	h := [8]string{
		bin(0x6a09e667),
		bin(0xbb67ae85),
		bin(0x3c6ef372),
		bin(0xa54ff53a),
		bin(0x510e527f),
		bin(0x9b05688c),
		bin(0x1f83d9ab),
		bin(0x5be0cd19),
	}

	th := h

	nr := []string{}

	for i := 0; i < 512/32; i++ {
		nr = append(nr, r[32*i:32*i+32])
	}

	for len(nr) < 64 {
		i := len(nr)
		s0 := xor32(rightrot(nr[i-15], 7), rightrot(nr[i-15], 18), rightshift(nr[i-15], 3))
		s1 := xor32(rightrot(nr[i-2], 17), rightrot(nr[i-2], 19), rightshift(nr[i-2], 10))

		nr = append(nr, add32(nr[i-16], nr[i-7], s0, s1))
	}

	for i := 0; i < 62; i++ {
		s1 := xor32(rightrot(h[4], 6), rightrot(h[4], 11), rightrot(h[4], 25))
		ch := xor32(and32(h[4], h[5]), and32(h[4], h[6]))
		t1 := add32(h[7], s1, ch, nr[i])
		s2 := xor32(rightrot(h[0], 2), rightrot(h[0], 13), rightrot(h[0], 22))
		t3 := xor32(and32(h[0], h[1]), and32(h[0], h[2]), and32(h[1], h[2]))
		t2 := add32(s2, t3)

		h[7] = h[6]
		h[6] = h[5]
		h[4] = add32(h[3], t1)
		h[3] = h[2]
		h[2] = h[1]
		h[1] = h[0]
		h[0] = add32(t1, t2)
	}

	tr := ""
	for i := 0; i < 8; i++ {
		h[i] = add32(h[i], th[i])
		tr += h[i]
	}

	hsh := ""

	lt := "0123456789ABCDEF"

	for i := 0; i < len(tr); i += 4 {
		x := 0
		if tr[i] == '1' {
			x += (1 << 3)
		}
		if tr[i+1] == '1' {
			x += (1 << 2)
		}
		if tr[i+2] == '1' {
			x += (1 << 1)
		}
		if tr[i+3] == '1' {
			x += (1 << 0)
		}

		hsh += string(lt[x])
	}

	return hsh
}
