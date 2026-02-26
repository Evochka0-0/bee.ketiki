import { Home, showNotification, DeleteProduct, bouquet_images } from "./utils.js";

let summ = 0;

async function loadCart(cart_data) {
  summ = 0;

  /** @type {Array<{id: number, count: number}>} */
  let products_info = JSON.parse(cart_data);

  let ids = [];
  products_info.forEach((element) => {
    ids.push(element.id_bouquet);
  });

  //передаем в fetch список id, получаем информацию для этих букетов в массиве строк
  if (ids.length > 0) {
    const products_container = document.getElementById("products_container");
    const loading_patch = document.getElementById("loading_patch");
    products_container.removeChild(loading_patch);

    return fetch("/cart", {
      //return чтобы await понимал, чего ему ждать
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ ids: ids }),
    })
      .then((response) => {
        return response.json();
      })
      .then((bouquets) => {
        bouquets.forEach((bouquet) => {
          const product_card = document.createElement("div");
          product_card.className = "product_card";

          // проверяем тип товара и выбираем формат картинки
          const image_bouquet = document.createElement("img");
          image_bouquet.className = "image_bouquet";
          bouquet_images(bouquet, image_bouquet);
        
          const product_info = document.createElement("div");
          product_info.className = "product_info";
          const product_name = document.createElement("h2");
          product_name.className = "product_name";
          product_name.textContent = bouquet.name;
          product_info.appendChild(product_name);

          const product_quantity = document.createElement("div");
          product_quantity.className = "product_quantity";
          const count_product = document.createElement("span");
          count_product.className = "count_product";

          let count = 1;
          products_info.forEach((element) => {
            if (element.id_bouquet == bouquet.id_bouquet) {
              count = element.count;
            }
          });
          count_product.textContent = count + " шт";
          product_quantity.appendChild(count_product);

          const product_price = document.createElement("div");
          product_price.className = "product_price";
          const summ_price = document.createElement("span");
          summ_price.className = "summ_price";
          let prise_product = bouquet.price * count;
          summ += prise_product;
          summ_price.textContent = prise_product.toLocaleString("ru-RU") + " ₽";
          product_price.appendChild(summ_price);

          const del_icon = document.createElement("img");
          del_icon.id = "del_icon";
          del_icon.src = "/images/del_icon.png";
          del_icon.addEventListener("click", () => {
            DeleteProduct(bouquet.id_bouquet);
            window.location.reload();
          });

          product_card.appendChild(image_bouquet);
          product_card.appendChild(product_info);
          product_card.appendChild(product_quantity);
          product_card.appendChild(product_price);
          product_card.appendChild(del_icon);

          products_container.appendChild(product_card);
        });

        const final_count = document.getElementById("final_count");
        let summ_count = 0;
        products_info.forEach((element) => {
          summ_count += element.count;
        });
        final_count.textContent = summ_count + " шт.";

        const final_summ = document.getElementById("final_summ");
        final_summ.textContent = summ.toLocaleString("ru-RU") + " ₽";
      });
  }
}

async function initiatePayment(order_id, amount) {
  const response = await fetch("/payments", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ id_order: order_id, amount: amount }),
  });
  if (!response.ok) {
    throw new Error("Ошибка инициирования оплаты: " + response.status);
  }
  const data = await response.json();
  if (data.status !== "pending") {
    throw new Error("Оплата уже обработана");
  }
  return data.payment_ref;
}

async function confirmPayment(order_id, payment_ref) {
  const response = await fetch("/payments/confirm", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ id_order: order_id, payment_ref: payment_ref }),
  });
  if (!response.ok) {
    throw new Error("Ошибка подтверждения оплаты: " + response.status);
  }
  const data = await response.json();
  if (data.status !== "paid") {
    throw new Error("Оплата не подтверждена");
  }
  return data;
}

function CreateOrder(cart_data, client_id) {
  const pay_button = document.getElementById("pay_button");
  /** @type {Array<{id_bouquet: number, count: number}>} */
  let products_info;

  const data_time = document.getElementById('datetime-local_deadline');
  pay_button.addEventListener("click", async () => {

    // дедлайн
    const deadline_value = data_time.value;
    //проверять на null не надо, потому что стоит значение по умолчанию
    products_info = JSON.parse(cart_data);
    const orderData = {
      client_id: client_id,
      total_cost: summ,
      deadline: deadline_value,
      items: products_info,
    };

    try {
      // Создаём заказ
      const orderResponse = await fetch("/orders", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(orderData),
      });
      if (!orderResponse.ok) {
        throw new Error("Ошибка создания заказа: " + orderResponse.status);
      }
      const orderDataResponse = await orderResponse.json();
      const order_id = orderDataResponse.id_order;

      // Инициируем оплату
      const payment_ref = await initiatePayment(order_id, summ);

      window.location.replace(`/payment.html?id_order=${order_id}&amount=${summ}&payment_ref=${payment_ref}`);
    } catch (error) {
      console.log("Произошла ошибка", error.message, error);
      showNotification(error.message);
    }
  });
}

async function loaddataTime(){
  const datatime = document.getElementById('datetime-local_deadline');
  const now = new Date();
  now.setHours(now.getHours() + 1);// минимум через час можно забрать

  const offset = now.getTimezoneOffset() * 60000; // перевод в милисек

  const min_input = new Date(now - offset).toISOString().slice(0, 16); // в локальном часовоп поясе

  datatime.min = min_input;
  datatime.value = min_input;

}

async function GetIdClient() {
  let client_id;
  try {
    const response = await fetch("/auth_access");
    if (!response.ok) {
      throw new Error("Не авторизован");
    }
    const data = await response.json();

    client_id = data.user_id;
  } catch (error) {
    console.error(error);
    console.log("проверка авторизации при получении id пользователя");
    window.location.href = "registration.html";
  }

  let cart_data = localStorage.getItem("cart");

  const cart_main = document.getElementById("cart_main");
  const left = document.getElementById("left");
  const right = document.getElementById("right");
  const emty_cart = document.getElementById("emty_cart");
  const btn_catalog = document.getElementById("btn_catalog");

  if (cart_data == null || JSON.parse(cart_data).length == 0) {
    cart_main.removeChild(left);
    cart_main.removeChild(right);

    //корзина пуста
    emty_cart.classList.remove("cart_empty_hidden"); //показываем текст В корзине пока пусто
    btn_catalog.addEventListener('click', () => {
      window.location.replace('index.html');
    })

    cart_main.classList.add("emty_main");
    return;
  }
  //читаем локал сторэч сейчас, чтобы заказ оформлялся по тем же данным, котрые отобразились на странице
  // чтобы у loadCart и CreateOrder были одинаковые данные
  await loadCart(cart_data);
  await loaddataTime();
  Home();
  CreateOrder(cart_data, client_id);
}

document.addEventListener("DOMContentLoaded", () => {
  GetIdClient();
  Home();
});
