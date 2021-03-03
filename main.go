package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Sprite struct {
	left  *ebiten.Image
	right *ebiten.Image
	up    *ebiten.Image
	down  *ebiten.Image
}

type Player struct {
	Sprite
	face     *ebiten.Image
	position ebiten.GeoM
}

type Game struct {
	PlayerOne  Player
	PlayerTwo  Player
	Conn       *websocket.Conn
	Num        string
	State      int
	Message    string
	MsgHistory string
}

func (g *Game) MovePlayerOne() {
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		g.PlayerOne.position.Translate(-1, 0)
		g.PlayerOne.face = g.PlayerOne.left
		g.Conn.WriteMessage(websocket.TextMessage, []byte(g.Num+"left"))
	} else
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.PlayerOne.position.Translate(0, 1)
		g.PlayerOne.face = g.PlayerOne.down
		g.Conn.WriteMessage(websocket.TextMessage, []byte(g.Num+"down"))
	} else
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		g.PlayerOne.position.Translate(0, -1)
		g.PlayerOne.face = g.PlayerOne.up
		g.Conn.WriteMessage(websocket.TextMessage, []byte(g.Num+"up"))
	} else
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		g.PlayerOne.position.Translate(1, 0)
		g.PlayerOne.face = g.PlayerOne.right
		g.Conn.WriteMessage(websocket.TextMessage, []byte(g.Num+"right"))
	} else
	if repeatingKeyPressed(ebiten.KeySemicolon) {
		g.State = 1
		g.Message = ""
	}
}

func (g *Game) Update() error {
	if g.State == 0 {
		g.MovePlayerOne()
	} else {
		g.Message += string(ebiten.InputChars())
		if repeatingKeyPressed(ebiten.KeyEnter) {
			g.State = 0
			err := g.Conn.WriteMessage(websocket.TextMessage, []byte("msg: "+g.Message))
			g.Message = ""
			if err != nil {
				fmt.Println(err)
			}
		}
		if repeatingKeyPressed(ebiten.KeyBackspace) {
			if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
				if len(g.Message) >= 1 {
					g.Message = g.Message[:len(g.Message)-1]
				}
			}
		}
	}
	return nil
}

func getImg(path string) *ebiten.Image {
	file, err := ebitenutil.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	out := ebiten.NewImageFromImage(img)
	return out
}

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	runnerImage *ebiten.Image
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Green)
	screen.DrawImage(g.PlayerTwo.face, &ebiten.DrawImageOptions{GeoM: g.PlayerTwo.position})
	screen.DrawImage(g.PlayerOne.face, &ebiten.DrawImageOptions{GeoM: g.PlayerOne.position})
	ebitenutil.DebugPrint(screen, g.MsgHistory)
	if g.Message != "" {
		ebitenutil.DebugPrintAt(screen, g.Message, 20, screenHeight-40)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func CreatePlayer(ws *websocket.Conn) (Player, Player, int, error) {
	ok := make(chan bool)
	start := make(chan bool)
	ready := make(chan bool)
	go func(k chan<- bool, q chan<- bool) {
		for {
			_, resp, err := ws.ReadMessage()
			fmt.Println("RESPONSE:",string(resp))
			if string(resp) == "ready" {
				fmt.Println("Ready to launch")
				q<- true
				return
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			if string(resp) == "taken" {
				fmt.Println("Pick another")
			}
			if string(resp) == "ok" {
				fmt.Println("Ok ok..")
				k <- true
			}
		}
	}(ok, start)
	var err error
	var choice int

	go func(k <-chan bool, d chan<- bool) {
		for {
			select {
			case <-k:
				d <- true
				return
			case <-time.After(time.Second):
				fmt.Println("Pick a player [1 or 2]")
				_, err = fmt.Scanf("%d", &choice)
				if err != nil {
					log.Fatal(err)
				}
				str := strconv.Itoa(choice)
				err = ws.WriteMessage(websocket.TextMessage, []byte("#choice:"+str))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}(ok, ready)

	<-ready

	fmt.Println("Yes yes")
	var player Player
	var friend Player
	if choice == 1 {
		player = Player{
			position: ebiten.GeoM{},
			Sprite: Sprite{
				left:  getImg("images/man_l.png"),
				right: getImg("images/man_r.png"),
				up:    getImg("images/man_u.png"),
				down:  getImg("images/man_d.png"),
			},
		}
		friend = Player{
			position: ebiten.GeoM{},
			Sprite: Sprite{
				left:  getImg("images/friend_l.png"),
				right: getImg("images/friend_r.png"),
				up:    getImg("images/friend_u.png"),
				down:  getImg("images/friend_d.png"),
			},
		}
		player.face = player.left
		friend.face = friend.right
		player.position.Translate(10, 10)
		friend.position.Translate(screenWidth-40, screenHeight-40)
	}
	if choice == 2 {
		player = Player{
			position: ebiten.GeoM{},
			Sprite: Sprite{
				left:  getImg("images/friend_l.png"),
				right: getImg("images/friend_r.png"),
				up:    getImg("images/friend_u.png"),
				down:  getImg("images/friend_d.png"),
			},
		}
		friend = Player{
			position: ebiten.GeoM{},
			Sprite: Sprite{
				left:  getImg("images/man_l.png"),
				right: getImg("images/man_r.png"),
				up:    getImg("images/man_u.png"),
				down:  getImg("images/man_d.png"),
			},
		}
		friend.face = friend.left
		player.face = player.right
		friend.position.Translate(10, 10)
		player.position.Translate(screenWidth-40, screenHeight-40)
	}

	fmt.Println("Waiting for other players...")
	<- start
	return player, friend, choice, err
}

type jwt struct {
	Token string
}

func getToken() (string, error) {
	link := "http://localhost:8056/login"
	//link := "https://tmp.mama.sh/api/login"
	var jsonStr = []byte(`{"username":"martin", "password":"T3stpass!"}`)
	req, err := http.NewRequest("POST", link, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var token jwt
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
func main() {
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	//TODO Will be used to close the connection later
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//link := "tmp.mama.sh"
	//u := url.URL{Scheme: "wss", Host: link, Path: "/api/game"}
	link := "localhost:8056"
	u := url.URL{Scheme: "ws", Host: link, Path: "/game"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("Authorization", "bearer "+token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	player, friend, choice, err := CreatePlayer(c)

	//SETUP GAME INTERACTION HAS ENDED HERE AND REMAINING SETUP WILL BE MADE BEFORE EXECUTION


	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Backend Game")
	game := &Game{
		PlayerOne: player,
		PlayerTwo: friend,
		Conn:      c,
		Num:       strconv.Itoa(choice),
	}

	go func() {
		var str string
		if choice == 1 {
			str = "2"
		} else {
			str = "1"
		}
		left := str + "left"
		right := str + "right"
		up := str + "up"
		down := str + "down"
		for {
			_, move, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println("RECV:", string(move))
			if string(move) == left {
				game.PlayerTwo.position.Translate(-1, 0)
				game.PlayerTwo.face = game.PlayerTwo.left
			} else
			if string(move) == down {
				game.PlayerTwo.position.Translate(0, 1)
				game.PlayerTwo.face = game.PlayerTwo.down
			} else
			if string(move) == up {
				game.PlayerTwo.position.Translate(0, -1)
				game.PlayerTwo.face = game.PlayerTwo.up
			} else
			if string(move) == right {
				game.PlayerTwo.position.Translate(1, 0)
				game.PlayerTwo.face = game.PlayerTwo.right
			} else
			if string(move[:3]) == "msg" {
				game.MsgHistory += "\n" + string(move[5:])
			}
		}
	}()


	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
