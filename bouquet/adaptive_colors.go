package bouquet

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"myproject/models"
	"myproject/utils"
	"net/http"
	"os"

	"github.com/cenkalti/dominantcolor"
	"github.com/lucasb-eyer/go-colorful"
)

func (b BouquetBase) BouquetColors(w http.ResponseWriter, r *http.Request) *utils.AppError {
	if r.Method != "GET" {
		return &utils.AppError{
			Err:     errors.New("MethodNotAllowed"),
			Message: "Метод не поддерживается",
			Code:    http.StatusMethodNotAllowed,
		}
	}
	db := b.BDB.DB
	context := r.Context()
	query := "SELECT * FROM `base_color_palette` WHERE 1"

	rows, err := db.QueryContext(context, query)
	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}

	palette := []models.BaseColorPalette{}
	defer rows.Close()
	for rows.Next() {
		color := models.BaseColorPalette{}
		err := rows.Scan(&color.IDBaseColor, &color.Hex, &color.Name)
		if err != nil {
			return &utils.AppError{
				Err:     err,
				Message: "Не удалось распаковать данные",
				Code:    http.StatusInternalServerError,
			}
		}
		palette = append(palette, color)
	}

	// запрос к букетам, проверяем что везде записаны цвета
	query = "SELECT id_bouquet, image_url FROM bouquets WHERE dominate_color = '' OR id_base_color = 1" // id = 1 значит что цвет не указан
	rows.Close()
	rows, err = db.QueryContext(context, query)
	if err != nil {
		message := fmt.Sprintf("Ошибка SQL: %s", err)
		return &utils.AppError{
			Err:     err,
			Message: message,
			Code:    http.StatusInternalServerError,
		}
	}

	type BouquetColorsCheck struct {
		IDBouquet int    `json:"id_bouquet"`
		ImageUrl  string `json:"image_url"`
	}
	datas := []BouquetColorsCheck{} // массив id букетов, в которых домин и базовый цвета не определены
	defer rows.Close()
	for rows.Next() {
		data := BouquetColorsCheck{}
		err = rows.Scan(&data.IDBouquet, &data.ImageUrl)
		if err != nil {
			message := fmt.Sprintf("Ошибка SQL: %s", err)
			return &utils.AppError{
				Err:     err,
				Message: message,
				Code:    http.StatusInternalServerError,
			}
		}
		datas = append(datas, data)
	}

	// если записи с цветом нет, вычисляем доминирующий и базовый цвет, записываем в таблицу
	fmt.Println("len(datas) ==", len(datas))
	if len(datas) == 0 {
		return &utils.AppError{
			Err:     errors.New("Цвета определены для всех букетов"),
			Message: "Цвета определены для всех букетов",
			Code:    http.StatusOK,
		}
	}

	for i := 0; i < len(datas); i++ {
		mess := fmt.Sprintf("ОБНОВЛЯЕМ datas[%d]///////////////////", i)
		fmt.Println(mess)
		dominate_color, app_err := GetDominateColor("X:/bee.ketiki — last/static" + datas[i].ImageUrl)
		if app_err != nil {
			return &utils.AppError{
				Err:     app_err,
				Message: app_err.Message,
				Code:    app_err.Code,
			}
		}
		// в бд записываем понятную строку в формате  #FF0000
		hsv_dominate_color := dominate_color.Hex()
		// сравниваем с цветами из палитры
		id_base_color := DefineBaseColor(dominate_color, palette)

		fmt.Println("id_base_color = ", id_base_color)

		query = "UPDATE bouquets SET dominate_color = ?, id_base_color = ? WHERE id_bouquet = ?"

		result, err := db.ExecContext(context, query, hsv_dominate_color, id_base_color, datas[i].IDBouquet)
		if err != nil {
			// Если ошибка SQL - логируем и идем дальше
			fmt.Printf("Ошибка обновления букета %d: %v\n", datas[i].IDBouquet, err)
			continue
		}

		rows, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("Ошибка получения RowsAffected для букета %d: %v\n", datas[i].IDBouquet, err)
			continue
		}

		if rows == 0 {
			// Это не критическая ошибка. Возможно, цвет не изменился.
			// Просто пишем в консоль и идем дальше.
			fmt.Printf("Букет %d: данные не изменились (0 строк обновлено)\n", datas[i].IDBouquet)
			continue
		}

		fmt.Printf("Букет %d успешно обновлен!\n", datas[i].IDBouquet)
	}
	return utils.Form_response(w, "Успешно определена и записана цветовая гамма букетов", http.StatusOK)
}

func GetDominateColor(path string) (*colorful.Color, *utils.AppError) {
	file, err := os.Open(path)
	if err != nil {
		return nil, &utils.AppError{
			Err:     err,
			Message: "Не удалось открыть файл",
			Code:    http.StatusInternalServerError,
		}
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, &utils.AppError{
			Err:     err,
			Message: "Ошибка image.Decode",
			Code:    http.StatusInternalServerError,
		}
	}

	color := dominantcolor.Find(image) // нашли доминирующий цвет
	// Структура: color.RGBA 4 поля: R (Red), G (Green), B (Blue), A (Alpha/Прозрачность)

	c, _ := colorful.MakeColor(color)
	H, S, V := c.Hsv()

	if S < 0.6 {
		S = 0.8
	}
	if V < 0.6 {
		V = 0.7
	}
	domin_color := colorful.Hsv(H, S, V)
	return &domin_color, nil
}

func DefineBaseColor(domicolor *colorful.Color, palette []models.BaseColorPalette) int {
	// находим расстояние между точкой в пространстве доминантого цвета и каждого из базовых цветов
	//минимальное расстояние - значит базовый цвет найден

	var id_base_color int = 0
	var min_distance float64 = 100000000

	for i := range palette {
		base_hsv, _ := colorful.Hex("#" + palette[i].Hex)
		distance := domicolor.DistanceLab(base_hsv)

		if float64(distance) < min_distance {
			min_distance = float64(distance)
			id_base_color = palette[i].IDBaseColor
		}
	}

	return id_base_color
}
