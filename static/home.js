import { bouquet_images, showNotification } from "./utils.js";
import { RemoveOrderConainer } from "./admin/admin_home.js";
import { orderStatuses } from "./utils.js";
import { fetchStatuses } from "./utils.js";
import { removeLoadings } from "./utils.js";
    
import { loadAllOrders } from "./admin/admin_home.js";

export let currentId;

export function exit() {
  const exit_icon = document.getElementById("log_out");
  exit_icon.addEventListener("click", function () {
    fetch("/auth_access", {
      method: "PUT",
    })
      .then((response) => {
        if (response.ok) {
          return response.text();
        } else {
          throw new Error("Статус: " + response.status);
        }
      })
      .then((text) => {
        console.log("Ответ сервера: ", text);
        localStorage.removeItem("cart");
        window.location.replace("index.html");
      })
      .catch((error) => {
        console.error("Произошла ошибка", error.message);
      });
  });
}

export function modalWindow() {
  const editIcon = document.getElementById("edit_icon");
  const modal_overlay = document.getElementById("modal_overlay");
  const close_button = document.getElementById("close_modal");
  editIcon.addEventListener("click", function () {
    modal_overlay.classList.remove("modal-hidden");
  });

  close_button.addEventListener("click", function () {
    modal_overlay.classList.add("modal-hidden");
  });
}

function editInfo(currentId) {
  const last_name_input = document.getElementById("lastname").value;
  const first_name_input = document.getElementById("firstname").value;
  const phone_input = document.getElementById("phone").value;
  const email_input = document.getElementById("email").value;

  const clientData = {
    id_client: parseInt(currentId),
    last_name: last_name_input,
    first_name: first_name_input,
    phone: phone_input,
    email: email_input,
  };

  const modal_overlay = document.getElementById("modal_overlay");
  fetch("/clients", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(clientData),
  })
    .then((response) => response.text())
    .then((responseText) => {
      modal_overlay.classList.add("modal-hidden");
      showNotification(responseText);
      loadClientInfo(currentId);
    });
}

export function loadClientInfo(currentId) {
  const name = document.getElementById("my_name");
  const phone = document.getElementById("client_phone");
  const email = document.getElementById("client_email");

  fetch(`/clients?id_client=${currentId}`)//GET
    .then((response) => response.json())
    .then((data) => {
      name.textContent = `${data.last_name} ${data.first_name}`;
      phone.textContent = data.phone;
      email.textContent = data.email;

      document.getElementById("lastname").value = data.last_name || "";
      document.getElementById("firstname").value = data.first_name || "";
      document.getElementById("phone").value = data.phone || "";
      document.getElementById("email").value = data.email || "";
    });
}

/**
 * для перехода на страницу заказа
 * @param {object} card - div
 */
export async function ListenerForCard(card, id) {
  card.addEventListener("click", () => {
    window.location.href = `order.html?id=${id}`;
  })
}

/**
 * для создания карточки заказа.
 * @param {object} order - Объект заказа с данными от сервера.
 * @returns {Promise<HTMLElement>} - Готовый HTML-элемент карточки заказа.
 */
async function createOrderCard(order) {
  const orderCard = document.createElement("div");
  orderCard.className = "order_container";
  const orderHeader = document.createElement("div");
  orderHeader.className = "order_header";
  const orderDate = document.createElement("span");
  const date = new Date(order.created_at);
  orderDate.textContent = date.toLocaleDateString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });

  const deadl_div = document.createElement('div');
  deadl_div.className = "deadl_div";

  const deadline = document.createElement('h3');
  deadline.className = "deadline";

  const deadline_data =  new Date(order.deadline);
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

  const imagesContainer = document.createElement("div");
  imagesContainer.className = "order_images";
  try {
    const items = await fetch(`/order_items?id_order=${order.id_order}`).then((res) =>
      res.json(),
    );
    if (items && items.length > 0) {
      items.slice(0, 3).forEach((item, index) => {
        const img = document.createElement("img");
        bouquet_images(item, img);
        img.alt = "Букет";
        img.className = "order_bouquet_image";
        img.style.zIndex = items.length - index;
        img.style.transform = `translateX(${index * 30}px) rotate(${
          index * 5 - 5
        }deg)`;
        imagesContainer.appendChild(img);
      });
    }
  } catch (e) {
    console.error("Не удалось загрузить состав заказа:", e);
  }

  const totalCost = document.createElement("p");
  totalCost.className = "order_price";
  totalCost.textContent = `${order.total_cost.toLocaleString("ru-RU")} ₽`;

  const statusesContainer = document.createElement("div");
  statusesContainer.className = "order_statuses";

  const statusBadge = document.createElement("div");
  statusBadge.className = "order_status_badge";
  statusBadge.textContent = orderStatuses[order.id_status] || "Неизвестно";
  statusesContainer.appendChild(statusBadge);

  const paymentBadge = document.createElement("div");
  paymentBadge.className = `payment_status_badge ${order.payment_status}`; // 'paid' или 'pending'
  paymentBadge.textContent = 
    order.payment_status === "paid" ? "Оплачен" : "Ожидает оплаты";
  statusesContainer.appendChild(paymentBadge);

  orderCard.appendChild(orderHeader);
  orderCard.appendChild(imagesContainer);
  orderCard.appendChild(totalCost);
  orderCard.appendChild(statusesContainer);

  ListenerForCard(orderCard, order.id_order);

  return orderCard;
}

export async function loadOrders(currentId) {
  const ready_orders = document.getElementById("ready_orders");
  const in_process_orders = document.getElementById("in_process_orders");
  const old_orders = document.getElementById("old_orders");

  ready_orders.innerHTML = "";
  in_process_orders.innerHTML = "";
  old_orders.innerHTML = "";

  try {
    const orders = await fetch(`/orders?user_id=${currentId}`).then((res) => {
      if (!res.ok) throw new Error("Ошибка сети");
      return res.json();
    });

    removeLoadings();

    const orders_empty = document.getElementById("orders_empty");
    const my_orders_container = document.getElementById("my_orders");

    if (!orders || orders.length === 0) {
      my_orders_container.classList.add("orders_section-hidden");
      orders_empty.classList.remove("empty_state-hidden");
      return;
    }

    my_orders_container.classList.remove("orders_section-hidden");
    orders_empty.classList.add("empty_state-hidden");

    for (const order of orders) {
      const orderCard = await createOrderCard(order);
      const statusName = orderStatuses[order.id_status];

      if (statusName === "Готов к выдаче") {
        ready_orders.appendChild(orderCard);
      } else if (statusName === "Новый" || statusName === "Собирается") {
        in_process_orders.appendChild(orderCard);
      } else {
        old_orders.appendChild(orderCard);
      }
    }

    document.getElementById("ready_section").style.display =
      ready_orders.childElementCount > 0 ? "block" : "none";
    document.getElementById("in_process_section").style.display =
      in_process_orders.childElementCount > 0 ? "block" : "none";
    document.getElementById("old_section").style.display =
      old_orders.childElementCount > 0 ? "block" : "none";
  } catch (error) {
    console.error("Ошибка загрузки заказов:", error);
    removeLoadings();
    showNotification("Не удалось загрузить заказы.");
  }
}

async function GetIdClient() {
  try {
    const response = await fetch("/auth_access");
    if (!response.ok) throw new Error("Не авторизован");
    const data = await response.json();
    currentId = data.user_id;
    //добавляем роль
    let role = data.role;

    await fetchStatuses()

    switch (role){
      case "user":
        loadOrders(currentId);
        loadClientInfo(currentId);
        modalWindow();
        exit();
        break;
      case "admin":
        RemoveOrderConainer();
        loadClientInfo(currentId);
        modalWindow();
        exit();
        loadAllOrders();
        break;
      default:
        showNotification("Не удалось поределить права доступа");
    }

    
  } catch (error) {
    console.log(error.message);
    window.location.href = "registration.html";
  }
}

document.addEventListener("DOMContentLoaded", () => {
  GetIdClient();
  const form = document.getElementById("edit_form");
  form.addEventListener("submit", function (event) {
    event.preventDefault();
    editInfo(currentId);
  });
});
