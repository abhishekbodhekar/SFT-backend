package main

type customers struct {
	Email     string `json:"email" gorm:"primaryKey"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Height    int    `json:"height"`
	Weight    int    `json:"weight"`
	Passsword string `json:"password"`
	Gender    string `json:"gender"`
}

type HttpResp struct {
	Error   string      `json:"error"`
	Result  bool        `json:"result"`
	Message interface{} `json:"message"`
}
