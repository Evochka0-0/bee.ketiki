import { arrowsListener, loadSlides } from "./slider.js";
import { addToCart, DeleteProduct, getCartData, bouquet_images } from "./utils.js";
import { Home, Cart, offerToLogIn} from "./utils.js";
import { checkAuthorisation } from "./utils.js";

/**
 * @param {object} bouquets - Изначальные букеты или отфильтрованные
 */
function updateCards(bouquets) {
  const conteiner = document.getElementById("bouquets-container");
  conteiner.innerHTML = "";

  
  const cartData = getCartData();
  
  bouquets.forEach((bouquet) => {

    const card = document.createElement("a");
    const id = bouquet.id_bouquet;

    const count = cartData.get(id) || 0;

    card.href = `product.html?id=${id}`;
    
    card.className = "card";

    const title = document.createElement("h3");
    title.textContent = bouquet.name;

    const image = document.createElement("img");
    bouquet_images(bouquet, image);
    image.alt = bouquet.name;

    const price = document.createElement("h1");
    price.id = "price";
    price.textContent = bouquet.price;

    card.appendChild(title);
    card.appendChild(image);
    card.appendChild(price);

    if (count == 0) {
      const to_basket_button = document.createElement("button");
      to_basket_button.className = "buy_buttons";
      to_basket_button.textContent = "Положить в корзину";
      to_basket_button.addEventListener("click", function (event) {
        event.preventDefault(); // Чтобы ссылка не сработала
        event.stopPropagation(); // Чтобы событие не дошло до родителя (карточки-ссылки)
        addToCart(id);
        
        updateCards(bouquets);
      });
      card.appendChild(to_basket_button);
    } else {
      const buttons = document.createElement("div");
      buttons.style.display = "flex";
      buttons.style.justifyContent = "space-between";
      buttons.style.alignItems = "center";
      buttons.style.gap = "10px";

      const add_button = document.createElement("button");
      const delete_button = document.createElement("button");
      add_button.className = "incart_button";
      delete_button.className = "incart_button";
      add_button.textContent = `В корзине ${count}`;
      add_button.addEventListener("click", function (event) {
        event.preventDefault(); // Чтобы ссылка не сработала
        event.stopPropagation(); // Чтобы событие не дошло до родителя (карточки-ссылки)
        addToCart(id);
        
        updateCards(bouquets);
      });
      delete_button.textContent = `🗑️`;
      delete_button.style.width = "fit-content";
      delete_button.addEventListener("click", function (event) {
        event.preventDefault(); 
        event.stopPropagation();
        DeleteProduct(id);
        
        updateCards(bouquets);
      });
      buttons.appendChild(add_button);
      buttons.appendChild(delete_button);
      card.appendChild(buttons);
    }
    conteiner.appendChild(card);
  });
}

function loadBouquets() {
  const param = "all";
  // здесь в строку запроса нужно добавлять парпаметры фильтрации, если они есть
  fetch(`/bouquets?type=${param}`) //отправляет запрос на сервер по адресу /bouquets
    .then((resposne) => {
      return resposne.json();
    })
    .then((bouquets) => {
      updateCards(bouquets);
    });
}

/** 
 *  функция которая проходится по элементам фильтрации html, 
 *формирует строку запроса с параметрами фильтра
 *@param {string} category категория по которой проверяем отмеченные фильтры
 *@param {URLSearchParams} finalParams строка с параметрами фильтрации
*/
function GetParamsToFilter(category, finalParams){
  //АААААААААААААА
  const checkedInputs = document.querySelectorAll(`input[name="${category}"]:checked`);
  const selectedIds = Array.from(checkedInputs).map(input => input.value);

  //превращаем массив выбранных id в строку 
  selectedIds.forEach((id) => {
    finalParams.append(category, id);
  })
}
/**
 * Добавляет визуальную отметку выбранным цветам
 */
function initColorCheckboxes() {
    const colorDivs = document.querySelectorAll('.colorChoosediv');
    colorDivs.forEach(div => {
        const checkbox = div.querySelector('input[type="checkbox"]');
        if (!checkbox) return;
        
        // Устанавливаем начальное состояние
        if (checkbox.checked) {
            div.classList.add('selected');
        }
        
        // Слушаем изменения
        checkbox.addEventListener('change', function() {
            if (this.checked) {
                div.classList.add('selected');
            } else {
                div.classList.remove('selected');
            }
        });
    });
}
/**
 * вызывается из слушателя кнопки
 */
async function applyFilters(){
  const params = new URLSearchParams({type: "usual"});
  GetParamsToFilter("occasion_names", params);
  GetParamsToFilter("flowers_names", params);
  GetParamsToFilter("color_names", params);

  // и еще проверить фильтр на цену
  const inputPrice = document.getElementById("inputPrice");
  // Если ползунок найден — берем цену, если нет — ставим 0 или null
  let price_value = inputPrice ? inputPrice.value : 0;  

  if (price_value != 0){
    params.set("price", price_value);
  }

  // теперь строка с параметрами готова, и ее можно засунуть в fetch
  try{
    const resp = await fetch(`/bouquets?${params.toString()}`);
    if (!resp.ok) throw new Error("Ошибка загрузки");
    
    const filteredBouquets = await resp.json();
    updateCards(filteredBouquets); // Твоя функция отрисовки
  }catch (error) {
      console.error("Проблема с фильтрацией:", error);
  }
}

/**
 * Захуячивает чекбоксы и лейблы в индивидуальные контейнеры
 * 
 * @param {{Array}} objects
 * @param {string} category
 * @param {number} count
 */
function BuildCheckbox(objects, category, count){
  const div = document.createElement('div');
  let idName = "ChekboxLabel" + category + count.toString();
  div.id = idName;

  const checkbox = document.createElement('input');
  checkbox.type = "checkbox";
  let idCheck = "Checkbox" + category + count.toString();
  checkbox.id = idCheck;
  checkbox.value = objects[0];
  checkbox.name = category;

  const label = document.createElement('label');
  label.htmlFor = idCheck;
  let labelId = "Label" + category + count.toString();
  label.id = labelId;
  label.textContent = objects[1];// здесь давно написан индекс 1 а не 2 ало гемини ты чо такая невнимательная невыспалась?

  div.appendChild(checkbox);
  div.appendChild(label);

  return div;
}

/**
 * Создает элементы фильтрации в соответствии с данными на базе
 */
function loadFiltersElements() {
  const filters_container = document.getElementById("filters_container");

  // очищаем контейнер перед добавлением
  filters_container.innerHTML = '';

  const requests = [
    fetch('/price_extremes'),
    fetch('/occasions'),
    fetch('/colors'),
    fetch('/flowers')
  ];

  Promise.all(requests)
    .then(resp => {
      return Promise.all(resp.map(r => r.json()));
    })
    .then(([extremes, occasions, colors, flowers]) => {

      // ========== ПОЛЗУНОК ДЛЯ ВЫБОРА ЦЕНЫ ==========
      const priceContainer = document.createElement('div');
      priceContainer.className = 'price-filter-container';
      
      // заголовок для цены
      const priceTitle = document.createElement('div');
      priceTitle.className = 'filter-title';
      priceTitle.textContent = '💰 Цена до';
      priceContainer.appendChild(priceTitle);
      
      // контент для цены
      const priceContent = document.createElement('div');
      priceContent.className = 'filter-content';
      
      const inputPrice = document.createElement('input');
      inputPrice.id = "inputPrice";
      inputPrice.type = "range";
      inputPrice.max = extremes.maxPrice;
      inputPrice.min = extremes.minPrice;
      inputPrice.value = extremes.maxPrice;
      inputPrice.style.width = '100%';
      priceContent.appendChild(inputPrice);
      
      const priceValue = document.createElement('span');
      priceValue.className = 'price-value';
      priceValue.textContent = `${Math.round(extremes.maxPrice).toLocaleString('ru-RU')} ₽`;
      priceContent.appendChild(priceValue);
      
      inputPrice.addEventListener('input', function() {
        priceValue.textContent = `${Math.round(this.value).toLocaleString('ru-RU')} ₽`;
      });
      
      priceContainer.appendChild(priceContent);
      filters_container.appendChild(priceContainer);

      // ========== НАЗНАЧЕНИЕ ==========
      const occasions_div = document.createElement('div');
      occasions_div.id = "occasions_div";
      
      const occasionsTitle = document.createElement('div');
      occasionsTitle.className = 'filter-title';
      occasionsTitle.textContent = '✨ По назначению';
      occasions_div.appendChild(occasionsTitle);
      
      const occasionsContent = document.createElement('div');
      occasionsContent.className = 'filter-content';
      
      let occasionsCount = 0;
      occasions.forEach((occasion) => {
        let category = "occasion_names";
        const data = [
          occasion.occasion_name,
          occasion.occasion_name
        ];
        const div = BuildCheckbox(data, category, occasionsCount);
        occasionsContent.appendChild(div);
        occasionsCount++;
      })
      occasions_div.appendChild(occasionsContent);

      // ========== ЦВЕТЫ В БУКЕТЕ ==========
      const flowers_div = document.createElement('div');
      flowers_div.id = "flowers_div";
      
      const flowersTitle = document.createElement('div');
      flowersTitle.className = 'filter-title';
      flowersTitle.textContent = '🌸 Цветы в букете';
      flowers_div.appendChild(flowersTitle);
      
      const flowersContent = document.createElement('div');
      flowersContent.className = 'filter-content';
      
      let flowersCount = 0;
      flowers.forEach((flower) => {
        let category = "flowers_names";
        const data = [
          flower.name_flower,
          flower.name_flower
        ];
        const div = BuildCheckbox(data, category, flowersCount);
        flowersContent.appendChild(div);
        flowersCount++;
      })
      flowers_div.appendChild(flowersContent);

      // ========== ПАЛИТРА ==========
      const palitra_nahui = document.createElement('div');
      palitra_nahui.id = "palitra_nahui";
      
      const palitraTitle = document.createElement('div');
      palitraTitle.className = 'filter-title';
      palitraTitle.textContent = '🎨 Цветовая гамма';
      palitra_nahui.appendChild(palitraTitle);
      
      const palitraContent = document.createElement('div');
      palitraContent.className = 'filter-content';
      
      let colorCount = 0;
      colors.forEach((color) => {
        if (color.name === "не указан") return;
        
        const colorChoosediv = document.createElement('div');
        colorChoosediv.className = "colorChoosediv";

        let category = "color_names";
        let data = [
          color.name,
          color.name
        ];
        const div = BuildCheckbox(data, category, colorCount);

        const square = document.createElement('div');
        square.style.width = "25px";
        square.style.height = "25px";
        square.style.borderRadius = "50%";
        square.style.background = `#${color.hex}`;
        square.style.border = "2px solid white";
        square.style.boxShadow = "0 2px 6px rgba(0,0,0,0.2)";

        colorChoosediv.appendChild(square);
        colorChoosediv.appendChild(div);
        
        const label = colorChoosediv.querySelector('label');
        if (label) label.style.display = 'none';

        palitraContent.appendChild(colorChoosediv);
        colorCount++;
      })
      palitra_nahui.appendChild(palitraContent);
      
      // ========== КОНТЕЙНЕР ДЛЯ ТРЁХ БЛОКОВ В РЯД ==========
      const rowContainer = document.createElement('div');
      rowContainer.className = 'filters-row';
      
      rowContainer.appendChild(occasions_div);
      rowContainer.appendChild(flowers_div);
      rowContainer.appendChild(palitra_nahui);
      
      filters_container.appendChild(rowContainer);
      
      // ========== КНОПКА ==========
      const button = document.createElement('button');
      button.id = "filterButton";
      button.textContent = "Применить";
      button.addEventListener('click', applyFilters);
      filters_container.appendChild(button);
      
      initColorCheckboxes();
    })
    .catch(err => console.error("Один из запросов упал", err));
}

document.addEventListener("DOMContentLoaded", () => {
  //ждет когда страница будет загружена
  offerToLogIn();
  loadSlides();
  arrowsListener();
  loadFiltersElements();
  loadBouquets();
  Home();
  Cart();
});
