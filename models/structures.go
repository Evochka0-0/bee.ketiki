package models

type Bouquet struct {
	IDBouquet       int     `json:"id_bouquet"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	ImageUrl        string  `json:"image_url"`
	ReserveImageUrl string  `json:"reserve_image_url"`
	DominateColor   string  `json:"dominate_color"`
	IDBaseColor     int     `json:"id_base_color"`
	Type            string  `json:"type"`
}

type Client struct {
	IDClient   int    `json:"id_client"`
	Last_name  string `json:"last_name"`
	First_name string `json:"first_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
}

type LogInData struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Order struct {
	IDOrder         int     `json:"id_order"`
	ClientID        int     `json:"id_client"`
	StatusID        int     `json:"id_status"`
	TotalCost       float64 `json:"total_cost"`
	PaymentStatus   string  `json:"payment_status"`
	PaymentRef      string  `json:"payment_ref"`
	Deadline        string  `json:"deadline"`        // дата когда клиент хочет получить товар
	PickUp_DateTime string  `json:"pickup_datetime"` // даиа когда уже запбрали заказ !! ТОЛЬКО у выполненных заказов
	CreatedAt       string  `json:"created_at"`
}

type OrderStatus struct {
	IDStatus   int    `json:"id_status"`
	NameStatus string `json:"name"`
}

type OrderItemWithImage struct {
	Quantity        int    `json:"quantity"`
	ImageUrl        string `json:"image_url"`
	ReserveImageUrl string `json:"reserve_image_url"`
	Type            string `json:"type"`
}

type Session struct {
	IDSession int    `json:"id_session"`
	IDClient  int    `json:"id_client"`
	Token     string `json:"token"`
	ExpiersAt string `json:"expires_at"`
}

type Reviews struct {
	Last_name  string `json:"last_name"`
	First_name string `json:"first_name"`
	Message    string `json:"message"`
	Grade      int    `json:"grade"`
	CreatedAt  string `json:"created_at"`
}

type BaseColorPalette struct {
	IDBaseColor int    `json:"id_base_color"`
	Hex         string `json:"hex"`
	Name        string `json:"name"`
}
