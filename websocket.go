package main

import (
)

/*
type user struct {
	username string
	password string
}

type jwt struct {
	Token string
}

func getToken() (string, error){
	link := "https://tmp.mama.sh/api/login"
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

func ws(token string){
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	link := "tmp.mama.sh"
	u := url.URL{Scheme: "wss", Host: link, Path: "/api/game"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("Authorization", "bearer "+token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte("hejsa med digsa"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}

func main() {
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
	ws(token)
}

 */