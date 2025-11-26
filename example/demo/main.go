package main

import (
	"fmt"
	"strings"

	"github.com/realjustice/swiss-pairing-engine/src/gotha"
)

func main() {
	fmt.Println("=== 瑞士制编排测试程序 ===")

	// 初始化引擎和比赛
	g := gotha.NewGotha()
	tournament := gotha.NewTournament()

	// 准备10个选手，设置不同的等级分
	players := []struct {
		Name      string
		FirstName string
		Rating    int
	}{
		{"Zhang", "San", 2600},  // 选手1 - 高手
		{"Li", "Si", 2400},      // 选手2 - 较强
		{"Wang", "Wu", 2200},    // 选手3 - 中上
		{"Zhao", "Liu", 2000},   // 选手4 - 中等
		{"Chen", "Qi", 2000},    // 选手5 - 中等
		{"Liu", "Ba", 1800},     // 选手6 - 中下
		{"Yang", "Jiu", 1800},   // 选手7 - 中下
		{"Huang", "Shi", 1600},  // 选手8 - 较弱
		{"Zhou", "Shiyi", 1400}, // 选手9 - 新手
		{"Wu", "Shier", 1200},   // 选手10 - 入门
	}

	fmt.Println("📋 选手信息：")
	for i, p := range players {
		player := gotha.NewPlayer()
		player.SetName(p.Name)
		player.SetFirstName(p.FirstName)
		player.SetParticipatingStr("1111111111") // 参加10轮比赛
		player.SetRank(p.Rating)
		tournament.AddPlayer(player)
		fmt.Printf("  %2d. %s %s (Rating: %d)\n", i+1, p.Name, p.FirstName, p.Rating)
	}

	// 设置比赛参数
	gps := tournament.GetTournamentSet().GetGeneralParameterSet()
	gps.SetNumberOfRounds(5) // 进行5轮比赛
	g.SetTournament(tournament)
	g.SelectSystem("SWISS")
	t := g.GetTournament()

	// 进行多轮编排测试
	fmt.Println("\n" + strings.Repeat("=", 70))

	// 第一轮编排
	fmt.Println("\n🎯 第1轮编排：")
	round1Games := make([]*gotha.Game, 0)
	t.Pair(1).Walk(func(game *gotha.Game) (isStop bool) {
		fmt.Printf("  台号%d: 白方 %s %s  VS  黑方 %s %s\n",
			game.TableNumber,
			game.GetWhitePlayer().Name, game.GetWhitePlayer().FirstName,
			game.GetBlackPlayer().Name, game.GetBlackPlayer().FirstName)
		round1Games = append(round1Games, game)
		return false
	})

	// 设置第一轮比赛结果
	fmt.Println("\n📝 设置第1轮比赛结果...")
	for i, game := range round1Games {
		// 前三台白胜，后两台平局
		if i < 3 {
			t.SetGameResult(1, game.TableNumber, "RESULT_WHITEWINS")
			fmt.Printf("  台号%d: 白方胜\n", game.TableNumber)
		} else {
			t.SetGameResult(1, game.TableNumber, "RESULT_EQUAL")
			fmt.Printf("  台号%d: 和棋\n", game.TableNumber)
		}
	}

	// 第二轮编排
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("\n🎯 第2轮编排（根据第1轮成绩）：")

	// 显示选手当前积分
	fmt.Println("\n当前积分榜：")
	allPlayers := t.GetPlayers()

	// 需要先计算积分
	t.Pair(2).Walk(func(game *gotha.Game) (isStop bool) {
		fmt.Printf("  台号%d: 白方 %s %s  VS  黑方 %s %s\n",
			game.TableNumber,
			game.GetWhitePlayer().Name, game.GetWhitePlayer().FirstName,
			game.GetBlackPlayer().Name, game.GetBlackPlayer().FirstName)
		return false
	})

	// 第二轮比赛结果设置
	fmt.Println("\n📝 设置第2轮比赛结果...")
	round2Games := t.GetGamesFromRound(1) // roundNumber从0开始
	for i, game := range round2Games {
		if i%3 == 0 {
			t.SetGameResult(2, game.TableNumber, "RESULT_WHITEWINS")
			fmt.Printf("  台号%d: 白方胜\n", game.TableNumber)
		} else if i%3 == 1 {
			t.SetGameResult(2, game.TableNumber, "RESULT_BLACKWINS")
			fmt.Printf("  台号%d: 黑方胜\n", game.TableNumber)
		} else {
			t.SetGameResult(2, game.TableNumber, "RESULT_EQUAL")
			fmt.Printf("  台号%d: 和棋\n", game.TableNumber)
		}
	}

	// 显示所有选手信息
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("\n📊 选手详细信息：")
	for i, p := range allPlayers {
		fmt.Printf("  %2d. %s %s - KeyString: %s\n",
			i+1, p.Name, p.FirstName, p.GetKeyString())
	}

	// 分析分组信息
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("\n🔍 分析：")
	fmt.Println("  • 第1轮：初始编排，按等级分排序后折半配对（Split&Slip）")
	fmt.Println("  • 第2轮：全胜组 vs 全胜组，半分组 vs 半分组，零分组 vs 零分组")
	fmt.Println("\n  权重优先级：")
	fmt.Println("    1. 避免重复对阵 (权重: 5×10^14)")
	fmt.Println("    2. 黑白平衡     (权重: 10^6)")
	fmt.Println("    3. 最小化分差   (权重: 10^11)")
	fmt.Println("    4. 上下调优化   (权重: 10^8)")
	fmt.Println("    5. 种子优化     (权重: 10^9)")

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("\n✅ 测试完成！")

	// 显示统计信息
	fmt.Println("\n📈 统计信息：")
	fmt.Printf("  • 总选手数: %d\n", len(allPlayers))
	fmt.Printf("  • 已完成轮次: 2\n")
	fmt.Printf("  • 总对局数: %d\n", len(t.SortGameByTableNumber()))
}
