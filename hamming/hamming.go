package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

//Переворачивает строку
func Reverse(str string) (rev string) {
	for i := len(str) - 1; i >= 0; i-- {
		rev += string(str[i])
	}
	return rev
}

//Расчет контрольных битов(и вставка в слова)
func Count_CB(marsh []string, r map[int]string, n int, rev bool) {
	for k := range marsh {
		for ind, val := range r {
			pow := int(math.Pow(2, float64(ind))) - 1
			var s int
			for i := 0; i < n; i++ {
				s1, _ := strconv.Atoi(string(marsh[k][i]))
				s2, _ := strconv.Atoi(string(val[i]))
				s += s1 * s2
			}
			mod := strconv.Itoa(s % 2)
			marsh[k] = marsh[k][:pow] + mod + marsh[k][pow + 1:]
		}
		if rev {
			marsh[k] = Reverse(marsh[k])
		}
	}
}

//По контрольным битам находит ошибочный бит и при необходимости исправляет его
func Correct(marsh []string, lg int)  {
	for k := range marsh {
		var ind string
		for i := lg; i >= 0; i-- {
			pow := int(math.Pow(2, float64(i))) - 1
			ind += string(marsh[k][pow])
		}
		index, _ := strconv.ParseInt(ind, 2, lg + 2)
		if index != 0 {
			if marsh[k][index - 1] == '0' {
				marsh[k] = marsh[k][:index - 1] + "1" + marsh[k][index:]
			} else {
				marsh[k] = marsh[k][:index - 1] + "0" + marsh[k][index:]
			}
		}
	}
}

//Переводит октеты в восьмиразрядный двоичный код
func Transl_oct(okts []byte) (okt_bin []string) {
	for _, okt := range okts {
		bin := strconv.FormatInt(int64(okt), 2)
		if len(bin) < 8 {
			bin = strings.Repeat("0", 8 - len(bin)) + bin
		}
		okt_bin = append(okt_bin, bin)
	}
	return okt_bin
}

//Из 8 разрядов делает слова n разрядов
func Addition(okt_bin []string, w int, n int, rev bool) (marsh []string) {
	var i, k int
	ok := true
	for _, okt := range okt_bin {
		if ok {
			marsh = append(marsh, okt)
			ok = false
			i = w
		} else {
			if n - i > w {
				marsh[k] += okt
				i += w
			} else {
				marsh[k] += okt[:n - i]
				if rev {
					marsh[k] = Reverse(marsh[k])
				}
				marsh = append(marsh, okt[n - i:])
				k += 1
				i = w - (n - i)
			}
		}
		if len(marsh[k]) == n {
			if rev {
				marsh[k] = Reverse(marsh[k])
			}
			ok = true
		}
	}
	if len(marsh[k]) < n {
		marsh[k] += strings.Repeat("0", n - len(marsh[k]))
	}
	return marsh
}

//Составление матрицы преобразования
func Control_bits(lg int, n int) map[int]string {
	r := make(map[int]string)
	for i := 0; i <= lg; i++ {
		ok := true
		var str string
		pow := int(math.Pow(2, float64(i)))
		str += strings.Repeat("0", pow - 1)
		for j := 0; j < (n)/pow; j++ {
			if ok {
				str += strings.Repeat("1", pow)
			} else {
				str += strings.Repeat("0", pow)
			}
			ok = !ok
		}
		r[i] = str[:n]
	}
	return r
}

//Удаление контрольных битов из слова
func Remove_CB(marsh []string, lg int) {
	for k := range marsh {
		for i := 0; i <= lg; i++ {
			pow := int(math.Pow(2, float64(i))) - i - 1
			marsh[k] = marsh[k][:pow] + marsh[k][pow + 1:]
		}
	}
}

//Из n-разрядных слов составляет последовательность октетов, заданных в десятичном формате
func Decoder(marsh []string) (code []byte) {
	var all string
	for _, mar := range marsh {
		all += mar
	}
	var i int
	for i < len(all) {
		if i+8 < len(all) {
			byt, _ := strconv.ParseInt(all[i:i+8], 2, 10)
			code = append(code, byte(byt))
			i += 8
		} else {
			byt, _ := strconv.ParseInt(all[i:] + strings.Repeat("0", 8 - len(all[i:])), 2, 10)
			code = append(code, byte(byt))
			break
		}
	}
	for i := len(code) - 1; i >= 0; i-- {
		if code[i] == 0 {
			code = code[:i]
		} else {
			break
		}
	}
	return code
}

//Добавление контроьных битов = 0
func AddBits(marsh []string, lg int)  {
	for k := range marsh {
		for i := 0; i <= lg; i++ {
			pow := int(math.Pow(2, float64(i)))
			marsh[k] = marsh[k][:pow - 1] + "0" + marsh[k][pow - 1:]
		}
	}
}

//Функция декодирования
func Decoding(okts []byte, n int, lg int, r map[int]string, rev bool) string {
	okt_bin := Transl_oct(okts)
	//fmt.Println(okt_bin)

	marsh := Addition(okt_bin, 8, n, rev)
	//fmt.Println(marsh)

	Count_CB(marsh, r, n, !rev)
	//fmt.Println(marsh)

	Correct(marsh, lg)
	//fmt.Println(marsh)

	Remove_CB(marsh, lg)
	//fmt.Println(marsh)

	decode := Decoder(marsh)
	//fmt.Println(decode)

	decoder_inf := string(decode)
	//fmt.Println(decoder_inf)
	return decoder_inf
}

//Функция кодирования
func Coding(answer string, n int, lg int, r map[int]string, rev bool) []byte {
	answ_byte := []byte(answer)
	//fmt.Println(answ_byte)

	answ_bin := Transl_oct(answ_byte)
	//fmt.Println(answ_bin)

	answ_marsh := Addition(answ_bin, 8, n - lg - 1, rev)
	//fmt.Println(answ_marsh)

	AddBits(answ_marsh, lg)
	//fmt.Println(answ_marsh)

	Count_CB(answ_marsh, r, n, !rev)
	//fmt.Println(answ_marsh)

	code := Decoder(answ_marsh)
	return code
}

//Main
func main() {
	n := 34

	lg := int(math.Log2(float64(n)))

	r := Control_bits(lg, n)

	okts := []byte{112, 213, 42, 196, 128, 143, 152, 32, 92, 70, 85, 35, 200, 8, 88, 12, 178, 134, 96, 17, 104, 33, 153, 158, 196, 108, 2, 53, 71, 136, 1, 248, 4, 213}

	fmt.Println(Decoding(okts, n, lg, r, true))

	answer := "602335729944"

	check_dec := Coding(answer, n, lg, r, false)
	fmt.Println(check_dec)

	if Decoding(check_dec, n, lg, r, true) == answer {
		println("Okay, great!")
	} else {
		println("All bad(")
	}
}
