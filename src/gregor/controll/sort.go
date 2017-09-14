package controll

func getArrays(a []int) (a1, a2 []int, ret bool, b int) {
	var b1, b2 []int

	if len(a) == 1 {
		//证明最大的数出来了
		return b1, b2, true, a[0]
	}

	if len(a)%2 != 0 {

		for i := 0; i < len(a)-1; i = i + 2 {
			b1[i] = a[i]
			b2[i] = a[i+1]

		}
		return b1, b2, false, a[len(a)-1]
	}

	for i := 0; i < len(a); i += 2 {
		b1[i] = a[i]
		b2[i] = a[i+1]
	}

	return b1, b2, false, -1
}

func getBiger(a, b int) int {
	if a > b {
		return a
	}

	return b
}

//得出一轮中最大的数的数组
func findBigArray(a []int) ([]int, bool, int) {
	var c []int

	b1, b2, ret, b := getArrays(a)

	if ret == true {
		//找到该轮中的最大的数:b
		return c, true, b
	}

	for i := 0; i < len(b1); i++ {
		num := getBiger(b1[i], b2[i])
		c[i] = num
	}

	if b != -1 {
		//证明是奇数数组，将b留置下一个数组中
		c[len(b1)] = b
	}

	return c, false, 0
}

//func findBigerNums(a int) []int {
//	//从键盘输入一个数组a
//	var aSli, bSli []int

//	//基础数据asli
//	//保存每一轮的比较数据
//	b, ret, num := findBigArray(aSli)

//	c, ret2, num2 := findBigArray(b)

//	d, ret3, num3 := findBigArray(c)

//	e, ret4, num4 := findBigArray(d)

//	f, ret5, num5 := findBigArray(e)

//	g, ret6, num6 := findBigArray(f)

//	h, ret7, num7 := findBigArray(g)

//	i, ret8, num8 := findBigArray(h)

//	j, ret9, num9 := findBigArray(i)

//	k, ret10, num10 := findBigArray(j)

//	l, ret11, num11 := findBigArray(k)

//	m, ret12, num12 := findBigArray(l)

//	n, ret13, num13 := findBigArray(m)

//	o, ret14, num14 := findBigArray(n)

//	p, ret15, num15 := findBigArray(o)

//	q, ret16, num16 := findBigArray(p)

//	r, ret17, num17 := findBigArray(q)

//	s, ret18, num18 := findBigArray(r)

//	t, ret19, num19 := findBigArray(s)

//	//找到这
//	return aSli
//}
