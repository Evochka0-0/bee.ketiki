import { Home, Cart, showNotification, bouquet_images } from "./utils.js";
import { orderStatuses } from "./utils.js";
import { fetchStatuses } from "./utils.js";

/**
 * для создания карточки букета
 * @param {object} bouquet - Объект букета с данными от сервера.
 * @returns {Promise<HTMLElement>} - Готовый HTML-элемент карточки заказа.
 */
async function loadBouquetCard(bouquet){
  const card = document.createElement('a');
  card.className = "card";

  const what_how_much = document.createElement('div');
  what_how_much.className = "what_how_much";

  const name_b = document.createElement('h1');
  name_b.className = "name_b";
  name_b.textContent = bouquet.name;

  const quantity = document.createElement('h1');
  quantity.className = "quantity";
  quantity.textContent = `${bouquet.quantity} ${" шт"}`;

  what_how_much.appendChild(name_b);
  what_how_much.appendChild(quantity);

  const image = document.createElement('img');
  image.className = "image";
  bouquet_images(bouquet, image);

  const price = document.createElement('p');
  price.className = "price";
  price.textContent = bouquet.price;

  // описани
  
  const description = document.createElement("p");
  description.className = "description";
  description.textContent = bouquet.description;


  card.appendChild(what_how_much);
  card.appendChild(image);
  card.appendChild(description);
  card.appendChild(price);

  card.href = `product.html?id=${bouquet.id_bouquet}`;
  return card;
}

async function loadOrderInfo(){
  await fetchStatuses();
    const order_info = document.getElementById("order_info");
    const bouquets_container = document.getElementById("bouquets-container");
    const client_info = document.getElementById("client_info");

    let params = new URLSearchParams(document.location.search);
    const id = Number(params.get("id"));

  try{
    const order_data = await fetch(`/order_data/${id}`).then((response) =>
      response.json(),
    );
    // данные о клиенте
    const order_client = order_data.order_client;

    const name = document.createElement('p');
      name.id = "name";
      name.textContent = `${order_client.last_name} ${order_client.first_name}`;

      const phone = document.createElement('p');
      phone.id = "phone";
      phone.textContent = `${"Телефон: "} ${order_client.phone}`;

      const email = document.createElement('p');
      email.id = "email";
      email.textContent = `${"Электронная почта: "} ${order_client.email}`;

      client_info.appendChild(name);
      client_info.appendChild(phone);
      client_info.appendChild(email);


      // данные о заказе
      const orderHeader = document.createElement("div");
      orderHeader.className = "order_header";
      const orderDate = document.createElement("span");
      const date = new Date(order_client.created_at);
      orderDate.textContent = date.toLocaleDateString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      });

      const deadl_div = document.createElement('div');
      deadl_div.className = "deadl_div";

      const deadline = document.createElement('h3');
      deadline.className = "deadline";

      const deadline_data =  new Date(order_client.deadline);
      deadline.textContent = deadline_data.toLocaleDateString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });

      const deadline_mess = document.createElement('h3');
      deadline_mess.textContent = "Будет готов: ";

      deadl_div.appendChild(deadline_mess);
      deadl_div.appendChild(deadline);

      orderHeader.appendChild(orderDate);
      orderHeader.appendChild(deadl_div);

      order_info.appendChild(orderHeader);

      const summ_status = document.createElement('div');
      summ_status.id = "summ_status";
      const total_cost = document.createElement('h1');
      total_cost.id = "total_cost";
      total_cost.textContent = `${"Сумма: "} ${order_client.total_cost}`;

      const status = document.createElement('p');
      status.id = "status";
      status.textContent = orderStatuses[order_client.id_status];

      const payment_status = document.createElement('p');
      payment_status.id = "payment_status";
      payment_status.textContent = order_client.payment_status;

      summ_status.appendChild(total_cost);
      summ_status.appendChild(status);
      summ_status.appendChild(payment_status);

      order_info.appendChild(summ_status);

      //букеты в составе
      //перебирем массив букеиов

      const order_bouquets = order_data.order_bouquets;

      for (const bouquet of order_bouquets) {
        let card = await loadBouquetCard(bouquet);
        bouquets_container.appendChild(card);
      }
  }catch(err){
      showNotification("Не удалось загрузить данные, попробуйте позже");
      console.error(err);
  } 
}

document.addEventListener("DOMContentLoaded", () => {
  Home();
  Cart();
  loadOrderInfo();
});