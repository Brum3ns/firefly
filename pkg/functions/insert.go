package functions

/**Insert data payload, random-keyword*/
func InsertData(s, p string) string {
	s = PayloadInsert(s, p)

	//Random insert point of 'x' type (a,s,i):
	s = RandomInsert(s) // strings.ReplaceAll(s, G.Random, RandomCreate(G.M_Random))

	//fmt.Println("->", s)//DEBUG

	return s
}
