package main

//func main() {
//	f, err := tea.LogToFile("debug.log", "debug")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	ans1, ans2 := loginPrompt()
//	fmt.Println(ans1, ans2)
//
//	board := styled.NewBoard()
//	board.InitBaseMenu()
//	p := tea.NewProgram(board)
//	if _, err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	app := NewHttpApp()
	app.Run()
}
