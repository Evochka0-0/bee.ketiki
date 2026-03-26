package bouquet

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"myproject/models"
	"myproject/utils"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type BouquetBase struct {
	BDB *utils.Base // Embedding встраивание
}

type ListIdBouquets struct {
	IDs []int `json:"ids"`
}

func (b BouquetBase) ListBouquetsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	switch r.Method {
	case "POST":
		db := b.BDB.DB
		context := r.Context()

		var list_id ListIdBouquets

		err := json.NewDecoder(r.Body).Decode(&list_id)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}

		var query_elements []string
		var query_questions []string
		query_begin := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets WHERE id_bouquet IN ("
		query_elements = append(query_elements, query_begin)

		for i := 0; i < len(list_id.IDs); i++ {
			query_questions = append(query_questions, "?")
		}

		questions := strings.Join(query_questions, ",")

		query_elements = append(query_elements, questions)
		query_elements = append(query_elements, ")")

		query := strings.Join(query_elements, "")

		args := make([]any, len(list_id.IDs))

		for i := 0; i < len(list_id.IDs); i++ {
			args[i] = list_id.IDs[i]
		}

		rows, err := db.QueryContext(context, query, args...)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка SQL запроса",
				Code:    http.StatusInternalServerError,
			}
		}

		defer rows.Close()
		bouquets := []models.Bouquet{}

		for rows.Next() {
			bouquet := models.Bouquet{}
			err := rows.Scan(&bouquet.IDBouquet, &bouquet.Name, &bouquet.Description, &bouquet.Price, &bouquet.ImageUrl, &bouquet.ReserveImageUrl, &bouquet.IDBaseColor, &bouquet.Type)
			if err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Не удалось распаковать данные",
					Code:    http.StatusInternalServerError,
				}
			}
			bouquets = append(bouquets, bouquet)
		}

		return utils.Form_response(w, bouquets, http.StatusOK)
	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

func (b BouquetBase) BouquetIdHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	db := b.BDB.DB
	idStr := strings.TrimPrefix(r.URL.Path, "/bouquets/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &utils.AppError{
			Err:     err,
			Message: "Неверный ID",
			Code:    http.StatusBadRequest,
		}
	}
	switch r.Method {
	case "GET":
		data := models.Bouquet{}
		context := r.Context()
		query := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets WHERE id_bouquet = ?"
		err := db.QueryRowContext(context, query, id).Scan(&data.IDBouquet, &data.Name, &data.Description, &data.Price, &data.ImageUrl, &data.ReserveImageUrl, &data.IDBaseColor, &data.Type)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка отправки запроса SELECT WHERE id",
				Code:    http.StatusInternalServerError,
			}
		}
		return utils.Form_response(w, data, http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
}

//переписываем для фильтрации
// юудем использовать запрос
/*
SELECT b.`id_bouquet`, b.`name`, b.`description`, b.`price`, b.`image_url`, b.`reserve_image_url`, color_p.`hex`, color_p.`color_name`, b.`type`, o.`occasion_name`, GROUP_CONCAT(f.name_flower) as flowers_list
FROM `bouquets` b
INNER JOIN `base_color_palette` color_p ON  b.`id_base_color` = color_p.`id_base_color`
INNER JOIN `occasion` o  ON b.`id_occasion` = o.`id_occasion`
INNER JOIN `bouquet_structure` s ON b.`id_bouquet` = s.`id_bouquet`
INNER JOIN `flowers` f ON s.`id_flower` = f.`id_flower`
WHERE ...
GROUP BY b.`id_bouquet`
*/
func (b BouquetBase) BouquetsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	db := b.BDB.DB
	switch r.Method {
	case "GET":
		requ_type := r.URL.Query().Get("type")

		if requ_type != "usual" && requ_type != "special" && requ_type != "all" {
			return &utils.AppError{
				Err:     errors.New("unsupported type"),
				Message: "Несуществующий тип букета",
				Code:    http.StatusInternalServerError,
			}
		}

		//параметры для фильтрации
		colors_list := r.URL.Query()["color_names"]
		occasions_list := r.URL.Query()["occasion_names"]
		flowers_list := r.URL.Query()["flowers_names"]
		maxPrice := r.URL.Query().Get("price")
		// если пользователь отметил галочками несколько цветов,
		// например Ромашки, Розы, то ему надо выдавать букеты с Ромашками, букеты с Розами, а не букеты в которых обязатеольно И Ромашки И розы

		var rows *sql.Rows
		var err error
		context := r.Context()

		//query := "SELECT id_bouquet, name, description, price, image_url, reserve_image_url, id_base_color, type FROM bouquets"
		base_query := "SELECT b.id_bouquet, b.name, b.description, b.price, b.image_url, b.reserve_image_url, color_p.hex, color_p.color_name, b.type, o.occasion_name," +
			" GROUP_CONCAT(f.name_flower) as flowers_list" +
			" FROM bouquets b" +
			" LEFT JOIN base_color_palette color_p ON  b.id_base_color = color_p.id_base_color" +
			" LEFT JOIN occasion o  ON b.id_occasion = o.id_occasion" +
			" LEFT JOIN bouquet_structure s ON b.id_bouquet = s.id_bouquet" +
			" LEFT JOIN flowers f ON s.id_flower = f.id_flower"

		// добавляем where по параметрам
		/* WHERE ...*/
		var where_conditions []string
		var queryArgs []any

		if requ_type == "usual" || requ_type == "special" {
			whereTypestr := "b.type = ?"
			where_conditions = append(where_conditions, whereTypestr)
			queryArgs = append(queryArgs, requ_type)
		}

		if len(colors_list) > 0 {
			var colors []string
			colors_condition := "color_p.color_name IN ("
			for i := 0; i < len(colors_list); i++ {
				colors = append(colors, "?")
				queryArgs = append(queryArgs, colors_list[i])
			}
			fullColors := colors_condition + strings.Join(colors, ", ") + ")"
			where_conditions = append(where_conditions, fullColors)
		}

		if len(occasions_list) > 0 {
			var occasions []string
			occasions_condition := "o.occasion_name IN ("
			for i := 0; i < len(occasions_list); i++ {
				occasions = append(occasions, "?")
				queryArgs = append(queryArgs, occasions_list[i])
			}
			fullOccasions := occasions_condition + strings.Join(occasions, ", ") + ")"
			where_conditions = append(where_conditions, fullOccasions)

		}

		if len(flowers_list) > 0 {
			var flowers []string
			flowers_condition := "f.name_flower IN ("
			for i := 0; i < len(flowers_list); i++ {
				flowers = append(flowers, "?")
				queryArgs = append(queryArgs, flowers_list[i])
			}
			fullFlowers := flowers_condition + strings.Join(flowers, ", ") + ")"
			where_conditions = append(where_conditions, fullFlowers)
		}

		var price_condition string = ""
		if maxPrice != "" {
			price_condition = "b.price <= ?"
			queryArgs = append(queryArgs, maxPrice)
		}

		ending := " GROUP BY b.id_bouquet " // но если в массиве условий нет условия для фильтрации по цветам то эту часть вставлять не нужно
		fullQuery := base_query + ending
		// если массив условий не пустой, собиарем полностью
		if len(where_conditions) > 0 {
			whereQueryes := strings.Join(where_conditions, " OR ")
			if maxPrice != "" {
				whereQueryes += " AND " + price_condition
			}
			fullQuery = base_query + " WHERE " + whereQueryes + ending
		}
		rows, err = db.QueryContext(context, fullQuery, queryArgs...)

		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}

		defer rows.Close()
		bouquets := []models.BouquetForFilters{}
		for rows.Next() {
			bouquet := models.BouquetForFilters{}
			if err := rows.Scan(&bouquet.IDBouquet, &bouquet.Name, &bouquet.Description, &bouquet.Price, &bouquet.ImageUrl,
				&bouquet.ReserveImageUrl, &bouquet.ColorHex, &bouquet.ColorName, &bouquet.Type, &bouquet.OccasionName, &bouquet.FlowersList); err != nil {
				return &utils.AppError{
					Err:     err,
					Message: "Ошибка чтения данных о товарах",
					Code:    http.StatusInternalServerError,
				}
			}
			bouquets = append(bouquets, bouquet)
		}
		if err = rows.Err(); err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка итерации",
				Code:    http.StatusInternalServerError,
			}
		}

		return utils.Form_response(w, bouquets, http.StatusOK)

	case "POST":
		context := r.Context()

		bouquet_data := models.Bouquet{}

		err := json.NewDecoder(r.Body).Decode(&bouquet_data)

		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Некорректный формат данных",
				Code:    http.StatusBadRequest,
			}
		}
		query := "INSERT INTO bouquets(name, description, price, image_url) VALUES (?, ?, ?, ?)"

		result, err := db.ExecContext(context, query, bouquet_data.Name, bouquet_data.Description, bouquet_data.Price, bouquet_data.ImageUrl)

		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Ошибка RowsAffected()",
				Code:    http.StatusInternalServerError,
			}
		}
		if rows == 0 {
			return &utils.AppError{
				Err:     errors.New("ноль строк изменено"),
				Message: "Затронуто 0 строк",
				Code:    http.StatusInternalServerError,
			}
		}

		return utils.Form_response(w, "Новый товар добавлен", http.StatusOK)

	default:
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}

	}
}

func (b BouquetBase) MaxMinPricesBouquetsHandler(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != "GET" {
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}

	db := b.BDB.DB
	context := r.Context()
	query := "SELECT MIN(price) AS min_price, MAX(price) AS max_price FROM bouquets"

	var minPrice sql.NullFloat64
	var maxPrice sql.NullFloat64
	err := db.QueryRowContext(context, query).Scan(&minPrice, &maxPrice)

	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}

	// Извлекаем "чистое" число для JSON
	finalMinPrice := 0.0
	if minPrice.Valid {
		finalMinPrice = minPrice.Float64
	}

	finalMaxPrice := 0.0
	if maxPrice.Valid {
		finalMaxPrice = maxPrice.Float64
	}

	result := struct {
		MaxPrice float64 `json:"maxPrice"`
		MinPrice float64 `json:"minPrice"`
	}{
		MaxPrice: finalMaxPrice,
		MinPrice: finalMinPrice,
	}

	return utils.Form_response(w, result, http.StatusOK)
}
