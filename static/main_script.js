import { arrowsListener, loadSlides } from "./slider.js";
import { addToCart, DeleteProduct, getCartData, bouquet_images } from "./utils.js";
import { Home, Cart, offerToLogIn} from "./utils.js";
import { checkAuthorisation } from "./utils.js";

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
  fetch(`/bouquets?type=${param}`) //отправляет запрос на сервер по адресу /bouquets
    .then((resposne) => {
      return resposne.json();
    })
    .then((bouquets) => {
      updateCards(bouquets);
    });
}

document.addEventListener("DOMContentLoaded", () => {
  //ждет когда страница будет загружена
  offerToLogIn();
  loadSlides();
  arrowsListener();
  loadBouquets();
  Home();
  Cart();
});
