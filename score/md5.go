package score

import (
	"strings"
	"cqu/util"
)

const (
	School = "10611"
)

var hexcase = true
var chrsz = 8
var mode = 32

func hex_md5(A string) string {
	temp := str2binl(A)
	temp2 := core_md5(temp, int32(len(A)*chrsz))
	return binl2hex(temp2)
}

func str2binl(D string) []int32 {
	a := 2
	for B := 0; B < len(D)*chrsz; B += chrsz {
		a = (B >> 5) + 1
	}
	var C = make([]int32, a)
	A := (1 << uint(chrsz)) - 1
	for B := 0; B < len(D)*chrsz; B += chrsz {
		C[B>>5] = C[B>>5] | int32((int([]byte(charAt(D, B/chrsz))[0])&A)<<uint(B%32))
	}
	result := make([]int32, 0, 10)
	for _, c := range C {
		result = append(result, c)
	}
	return result
}

func core_md5(enter []int32, F int32) []int32 {
	temp := moveRight(F+64, 9)
	var a = (temp << 4) + 14 + 1
	K := make([]int32, a)
	for i := 0; i < len(K); i++ {
		K[i] = 0
	}
	for i := 0; i < len(enter); i++ {
		K[i] = enter[i]
	}
	K[F>>5] = K[F>>5] | (128 << uint(F%32))
	K[(moveRight(F+64, 9)<<4)+14] = F

	var J int32 = 1732584193
	var I int32 = -271733879
	var H int32 = -1732584194
	var G int32 = 271733878

	for C := 0; C < len(K); C += 16 {

		E := J
		D := I
		B := H
		A := G

		J = md5_ff(J, I, H, G, K[C+0], 7, -680876936)
		G = md5_ff(G, J, I, H, K[C+1], 12, -389564586)
		H = md5_ff(H, G, J, I, K[C+2], 17, 606105819)
		I = md5_ff(I, H, G, J, K[C+3], 22, -1044525330)
		J = md5_ff(J, I, H, G, K[C+4], 7, -176418897)
		G = md5_ff(G, J, I, H, K[C+5], 12, 1200080426)
		H = md5_ff(H, G, J, I, K[C+6], 17, -1473231341)
		I = md5_ff(I, H, G, J, K[C+7], 22, -45705983)
		J = md5_ff(J, I, H, G, K[C+8], 7, 1770035416)
		G = md5_ff(G, J, I, H, K[C+9], 12, -1958414417)
		H = md5_ff(H, G, J, I, K[C+10], 17, -42063)
		I = md5_ff(I, H, G, J, K[C+11], 22, -1990404162)
		J = md5_ff(J, I, H, G, K[C+12], 7, 1804603682)
		G = md5_ff(G, J, I, H, K[C+13], 12, -40341101)
		H = md5_ff(H, G, J, I, K[C+14], 17, -1502002290)
		// I = md5_ff(I, H, G, J, K[C + 15], 22, 1236535329)
		I = md5_ff(I, H, G, J, 0, 22, 1236535329)
		J = md5_gg(J, I, H, G, K[C+1], 5, -165796510)
		G = md5_gg(G, J, I, H, K[C+6], 9, -1069501632)
		H = md5_gg(H, G, J, I, K[C+11], 14, 643717713)
		I = md5_gg(I, H, G, J, K[C+0], 20, -373897302)
		J = md5_gg(J, I, H, G, K[C+5], 5, -701558691)
		G = md5_gg(G, J, I, H, K[C+10], 9, 38016083)
		// H = md5_gg(H, G, J, I, K[C + 15], 14, -660478335)
		H = md5_gg(H, G, J, I, 0, 14, -660478335)
		I = md5_gg(I, H, G, J, K[C+4], 20, -405537848)
		J = md5_gg(J, I, H, G, K[C+9], 5, 568446438)
		G = md5_gg(G, J, I, H, K[C+14], 9, -1019803690)
		H = md5_gg(H, G, J, I, K[C+3], 14, -187363961)
		I = md5_gg(I, H, G, J, K[C+8], 20, 1163531501)
		J = md5_gg(J, I, H, G, K[C+13], 5, -1444681467)
		G = md5_gg(G, J, I, H, K[C+2], 9, -51403784)
		H = md5_gg(H, G, J, I, K[C+7], 14, 1735328473)
		I = md5_gg(I, H, G, J, K[C+12], 20, -1926607734)
		J = md5_hh(J, I, H, G, K[C+5], 4, -378558)
		G = md5_hh(G, J, I, H, K[C+8], 11, -2022574463)
		H = md5_hh(H, G, J, I, K[C+11], 16, 1839030562)
		I = md5_hh(I, H, G, J, K[C+14], 23, -35309556)
		J = md5_hh(J, I, H, G, K[C+1], 4, -1530992060)
		G = md5_hh(G, J, I, H, K[C+4], 11, 1272893353)
		H = md5_hh(H, G, J, I, K[C+7], 16, -155497632)
		I = md5_hh(I, H, G, J, K[C+10], 23, -1094730640)
		J = md5_hh(J, I, H, G, K[C+13], 4, 681279174)
		G = md5_hh(G, J, I, H, K[C+0], 11, -358537222)
		H = md5_hh(H, G, J, I, K[C+3], 16, -722521979)
		I = md5_hh(I, H, G, J, K[C+6], 23, 76029189)
		J = md5_hh(J, I, H, G, K[C+9], 4, -640364487)
		G = md5_hh(G, J, I, H, K[C+12], 11, -421815835)
		// H = md5_hh(H, G, J, I, K[C + 15], 16, 530742520)
		H = md5_hh(H, G, J, I, 0, 16, 530742520)
		I = md5_hh(I, H, G, J, K[C+2], 23, -995338651)
		J = md5_ii(J, I, H, G, K[C+0], 6, -198630844)
		G = md5_ii(G, J, I, H, K[C+7], 10, 1126891415)
		H = md5_ii(H, G, J, I, K[C+14], 15, -1416354905)
		I = md5_ii(I, H, G, J, K[C+5], 21, -57434055)
		J = md5_ii(J, I, H, G, K[C+12], 6, 1700485571)
		G = md5_ii(G, J, I, H, K[C+3], 10, -1894986606)
		H = md5_ii(H, G, J, I, K[C+10], 15, -1051523)
		I = md5_ii(I, H, G, J, K[C+1], 21, -2054922799)
		J = md5_ii(J, I, H, G, K[C+8], 6, 1873313359)
		// G = md5_ii(G, J, I, H, K[C + 15], 10, -30611744)
		G = md5_ii(G, J, I, H, 0, 10, -30611744)
		H = md5_ii(H, G, J, I, K[C+6], 15, -1560198380)
		I = md5_ii(I, H, G, J, K[C+13], 21, 1309151649)
		J = md5_ii(J, I, H, G, K[C+4], 6, -145523070)
		G = md5_ii(G, J, I, H, K[C+11], 10, -1120210379)
		H = md5_ii(H, G, J, I, K[C+2], 15, 718787259)
		I = md5_ii(I, H, G, J, K[C+9], 21, -343485551)
		J = safe_add(J, E)
		I = safe_add(I, D)
		H = safe_add(H, B)
		G = safe_add(G, A)
	}
	result := make([]int32, 0, 4)
	if mode == 16 {
		result = append(result, I, H)
	} else {
		result = append(result, J, I, H, G)
	}
	return result
}

func safe_add(A, D int32) int32 {
	C := (A & 65535) + (D & 65535)
	B := (A >> 16) + (D >> 16) + (C >> 16)
	return (B << 16) | (C & 65535)
}

func moveRight(val int32, count int) int32 {
	var temp int32 = 0
	for i := 0; i < 32-count; i++ {
		temp = temp ^ (1 << uint(i))
	}
	val = val >> uint(count)
	val = val & temp
	return val
}

func bit_rol(A, B int32) int32 {
	return (A << uint(B)) | moveRight(A, int(32-B))
}

func md5_cmn(F, C, B, A, E, D int32) int32 {
	cf := safe_add(C, F)
	ad := safe_add(A, D)
	safeAdd := safe_add(cf, ad)
	bitRol := bit_rol(safeAdd, E)
	return safe_add(bitRol, B)
}
func md5_ff(C, B, G, F, A, E, D int32) int32 {
	return md5_cmn((B&G)|((^B)&F), C, B, A, E, D)
}
func md5_gg(C, B, G, F, A, E, D int32) int32 {
	return md5_cmn((B&F)|(G&(^F)), C, B, A, E, D)
}
func md5_hh(C, B, G, F, A, E, D int32) int32 {
	return md5_cmn(B^G^F, C, B, A, E, D)
}
func md5_ii(C, B, G, F, A, E, D int32) int32 {
	return md5_cmn(G^(B|(^F)), C, B, A, E, D)
}

func binl2hex(C []int32) string {
	B := ""
	if hexcase {
		B = "0123456789ABCDEF"
	} else {
		B = "0123456789abcdef"
	}
	D := ""
	for A := 0; A < len(C)*4; A++ {
		r := (C[A>>uint(2)] >> uint((A%4)*8+4)) & 15
		r2 := (C[A>>uint(2)] >> uint((A%4)*8)) & 15
		D += charAt(B, int(r)) + charAt(B, int(r2))
	}
	return D
}

func CquMd5Encrypted(id, pass string) string {
	passwordEncrypted := hex_md5(pass)
	addSchoolId := id + strings.ToUpper(util.SubString(passwordEncrypted, 0, 30)) + School
	encryptedWithUsername := strings.ToUpper(util.SubString(hex_md5(addSchoolId), 0, 30))
	return encryptedWithUsername
}

func charAt(str string, index int) string {
	bytes := []byte(str)
	lenght := len(bytes)
	if index < 0 {
		index = 0
	}
	if index > lenght {
		index = lenght
	}
	return string(bytes[index])
}
