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
function loadFiltersElements(){
  const filters_container = document.getElementById("filters_container");

  //находим максимальную и минимальную цену
  //список назначений букетов
  //узнаем, каких цветов есть букеты
  //узнаем какие цветки могут быть в составе

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
    // тут по идее нужно вызвать функцию которая отрисует каждый элемент либо сделать это сразу тут

    //_________ползунок для выбора цены
    const inputPrice = document.createElement('input');
    inputPrice.id = "inputPrice";
    inputPrice.type = "range";
    inputPrice.max = extremes.maxPrice;
    inputPrice.min = extremes.minPrice;
    inputPrice.value = extremes.maxPrice;
    //добавить рядом с ним маленький <span>, 
    // который будет обновляться и показывать текущее число, когда ползунок двигают (событие input)
    filters_container.appendChild(inputPrice);


    //_____СУКА ОКАКАШОНС 
    const occasions_div = document.createElement('div');
    occasions_div.id = "occasions_div";
    let count = 0;
    occasions.forEach((occasion) => {
      // вызов функции я того туда сюда манала, ее ж написать еще надо
      let category = "occasion_names";
      const data = [
        occasion.occasion_name,//value
        occasion.occasion_name//name
      ];
      const div = BuildCheckbox(data, category, count);
      occasions_div.appendChild(div);
      count ++;
    })


    //______цветы розы мимозы
    const flowers_div = document.createElement('div');
    flowers_div.id = "flowers_div";
    count = 0;
    flowers.forEach((flower) => {
      let category = "flowers_names";// поменять на flowers_list?
      const data = [
        flower.name_flower,
        flower.name_flower
      ];

      const div = BuildCheckbox(data, category, count);
      flowers_div.appendChild(div);
      count++;
    })

    //____ПАЛИТРА
    const palitra_nahui = document.createElement('div');
    palitra_nahui.id = "palitra_nahui";
    count = 0;
    colors.forEach((color) => {
      // создаем контейнер с цветным квадратиком
      const colorChoosediv = document.createElement('div');
      colorChoosediv.className = "colorChoosediv";

      let category = "color_names";
      let data = [
        color.name,
        color.name
      ] 
      const div = BuildCheckbox(data, category, count);

      const square = document.createElement('div');
      square.style.width = "10px";
      square.style.height  = "10px";
      square.style.background = `#${color.hex}`;

      colorChoosediv.appendChild(square);
      colorChoosediv.appendChild(div);

      palitra_nahui.appendChild(colorChoosediv);
      count ++;
    })
    
    const button = document.createElement('button');
    button.id = "filterButton";
    button.textContent = "Применить";
    button.addEventListener('click', applyFilters);

    filters_container.appendChild(occasions_div);
    filters_container.appendChild(flowers_div);
    filters_container.appendChild(palitra_nahui);
    filters_container.appendChild(button);
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
