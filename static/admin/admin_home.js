import { bouquet_images, orderStatuses } from "../utils.js";
import { removeLoadings } from "../utils.js";
import { showNotification } from "../utils.js";
import { ListenerForCard } from "../home.js";

export function RemoveOrderConainer(){
    const my_orders_container = document.getElementById("my_orders");
    my_orders_container.classList.add("orders_section-hidden");

    const admin_panel = document.getElementById("admin_panel");
    admin_panel.classList.remove("admin_panel_hidden");

    return;
}

// заполняем секции такие же как у пользователя, но отображаем ВСЕ ЗАКАЗЫ которые есть в системе
// GET orders 

async function createAdminOrderCard(order) {//изменять статус
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
  deadline_mess.textContent = "Дедлайн: ";

  deadl_div.appendChild(deadline_mess);
  deadl_div.appendChild(deadline);

  orderHeader.appendChild(orderDate);
  orderHeader.appendChild(deadl_div);

  //информация о клиенте в скрытом окошке
  const clientWrapper = document.createElement("div");
  clientWrapper.className = "dropdown";

  const title = document.createElement("span");
  title.textContent = "О клиенте ▾";
  title.className = "dropdown-title";

  const content = document.createElement("div");
  content.className = "dropdown-content hidden";

  const exit_icon = document.createElement("h3");
  exit_icon.className = "exit_icons";
  exit_icon.textContent = "❌";

  exit_icon.addEventListener('click', (e) => {
    e.stopPropagation(); 
      content.classList.toggle("hidden");
  })

  const name = document.createElement("h4");
  name.textContent = `Заказчик: ${order.last_name} ${order.first_name}`;
  const phone = document.createElement("h4");
  phone.textContent = `Телефон: ${order.phone}`;
  const email = document.createElement("h4");
  email.textContent = `Почта: ${order.email}`;

  content.appendChild(exit_icon);
  content.appendChild(name);
  content.appendChild(phone);
  content.appendChild(email);

  title.addEventListener('click', (e) => {
      e.stopPropagation();
      content.classList.toggle("hidden");
  });

  clientWrapper.appendChild(title);
  clientWrapper.appendChild(content);

  const imagesContainer = document.createElement("div");
  imagesContainer.className = "order_images";

  try{
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
  }catch(err){
    console.error("Не удалось загрузить состав заказа:", err);
  }

  const totalCost = document.createElement("p");
  totalCost.className = "order_price";
  totalCost.textContent = `${order.total_cost.toLocaleString("ru-RU")} ₽`;

  const statusesContainer = document.createElement("div");
  statusesContainer.className = "order_statuses";

  const selectorStatus = document.createElement("select");
  selectorStatus.className = "selector";

  // преобразуем обхект со статусами в удобный формат и перебираем  
  let statuses = Object.entries(orderStatuses);
  statuses.forEach(([id, name]) => {
    //добавляем в селектор статусы
    const option = document.createElement("option");
    option.value = id;
    option.textContent = name;
    if (Number(id) === order.id_status){
      option.selected = true;
    }
    selectorStatus.appendChild(option);
  })

  selectorStatus.addEventListener('click', function(event){
    event.preventDefault(); // Чтобы ссылка не сработала
    event.stopPropagation(); // Чтобы событие не дошло до родителя (карточки-ссылки)
  })

  selectorStatus.addEventListener('change', async(event) => {
    //когда админ нажал на статус, вычисляем его id, отправляем на сервер в PUT order с id_order
    const data = {
      id_status: Number(event.target.value),      
      id_order: order.id_order
    };

    try{
        const response = await fetch(`/admin/status`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      if (!response.ok){
        throw new Error("Не удалось обновить статус заказа", response.status);
      }
      else{
        const data =  await response.json();
        showNotification(data);
        window.location.reload();
      }
    }
    catch(error){
      console.error("Ошибка загрузки заказов:", error);
      removeLoadings();
      showNotification("Не удалось загрузить данные о заказах!");
    }
  })

  statusesContainer.appendChild(selectorStatus);

  const paymentBadge = document.createElement("div");
  paymentBadge.className = `payment_status_badge ${order.payment_status}`;
  paymentBadge.textContent =
    order.payment_status === "paid" ? "Оплачен" : "Ожидает оплаты";
  statusesContainer.appendChild(paymentBadge);

  orderCard.appendChild(orderHeader);
  orderCard.appendChild(clientWrapper);
  orderCard.appendChild(imagesContainer);
  orderCard.appendChild(totalCost);
  orderCard.appendChild(statusesContainer);

  ListenerForCard(orderCard, order.id_order);
  return orderCard;

}


export async function loadAllOrders() {
  const ready_orders_adm = document.getElementById("ready_orders_adm");
  const in_process_orders_adm = document.getElementById("in_process_orders_adm");

  // Очищаем контейнеры перед загрузкой
  ready_orders_adm.innerHTML = "";
  in_process_orders_adm.innerHTML = "";

  try {
    const orders = await fetch(`/admin/orders`).then((response) => {
        if (!response.ok) throw new Error("Ошибка загрузки заказов:", response.status);
        return response.json();
    });
    

    removeLoadings();

    const adm_orders_empty = document.getElementById("adm_orders_empty");
    const admin_panel = document.getElementById("admin_panel");

    if (!orders || orders.length === 0) {
      admin_panel.classList.add("admin_panel_hidden");
      adm_orders_empty.classList.remove("adm_empty_state-hidden");
      return;
    }

    adm_orders_empty.classList.add("adm_empty_state-hidden");

    for (const order of orders) {
        const orderCard = await createAdminOrderCard(order);
        const statusName = orderStatuses[order.id_status];

        if (statusName === "Готов к выдаче") {
            ready_orders_adm.appendChild(orderCard);
        } else if (statusName === "Новый" || statusName === "Собирается") {
            in_process_orders_adm.appendChild(orderCard);
        }
    }

    // Скрываем пустые секции
    document.getElementById("ready_adm").style.display =
      ready_orders_adm.childElementCount > 0 ? "block" : "none"; // 🙈🙈 че это нахуй???????
    document.getElementById("in_process_adm").style.display =
      in_process_orders_adm.childElementCount > 0 ? "block" : "none";// 🙈🙈 че это нахуй???????
    
  } catch (error) {
    console.error("Ошибка загрузки заказов:", error);
    removeLoadings();
    showNotification("Не удалось загрузить заказы.");
  }
}

// еще нужны кнопки для смены статуса заказа
// PUT order {id}