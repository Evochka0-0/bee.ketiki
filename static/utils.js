export function bouquet_images(bouquet, image_bouquet){
  if(bouquet.type != "usual"){
    image_bouquet.src = bouquet.reserve_image_url; // витринный
    }else{
      image_bouquet.src = bouquet.image_url; // обычный
    }
}

export function checkAuthorisation(path) {
  //принимает путь к странице, к которой нужен переход с проверкой авторизации
  fetch("/auth_access")
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        // Если статус 401 или другой плохой, "выбрасываем" ошибку\
        throw new Error("Пользователь не авторизован");
      }
    })
    .then((data) => {
      window.location.href = path;
    })
    .catch((error) => {
      console.log("Произошла ошибка" + error.message);
      window.location.href = "registration.html";
    });
}

export function offerToLogIn(){
  if (!document.cookie.includes("session_token")) {
    const div = document.createElement('div');
    div.className = "offer_to_login_div";
    
    const href_text = document.createElement('a');
    href_text.className = "href_text";
    href_text.textContent = "Зарегестрироваться/Войти";
    href_text.href = "registration.html";

    div.appendChild(href_text);
    const bouquets_container = document.getElementById("bouquets-container");
    document.querySelector("main").insertBefore(div, bouquets_container);
  }
}

export function Home() {
  const home = document.getElementById("home");
  home.addEventListener("click", function () {
    //проверка авторизации, если нет -> переход на registration.html
    checkAuthorisation("home.html");
  });
}

export function Cart() {
  const cart_icon = document.getElementById("cart_id");

  cart_icon.addEventListener("click", function () {
    checkAuthorisation("cart.html");
    //window.location.href = 'cart.html';
  });
}

export function showNotification(responseText) {
  const mes_container = document.createElement("div");
  mes_container.className = "toast";
  mes_container.textContent = responseText;

  document.querySelector("main").appendChild(mes_container);

  setTimeout(() => {
    mes_container.classList.add("toast-visible");
  }, 100);

  setTimeout(() => {
    mes_container.classList.remove("toast-visible");
    setTimeout(() => {
      mes_container.remove();
    }, 500);
  }, 3000);
}

export function DeleteProduct(id_todel) {
  if (!document.cookie.includes("session_token")) {
    showNotification(
      "Пожалуйста, зарегестрируйтесь/войдите, чтобы взаимодействовать с корзиной",
    );
  }

  //читаем актуальный local Storage

  let cart_data = localStorage.getItem("cart");

  /** @type {Array<{id: number, count: number}>} */
  let products_info = JSON.parse(cart_data);
  //копируем все кроме того который удаляем

  const new_product_info = products_info.filter((item) => item.id_bouquet !== id_todel);

  localStorage.setItem("cart", JSON.stringify(new_product_info));

}

export function addToCart(id) {
  if (!document.cookie.includes("session_token")) {
    showNotification(
      "Пожалуйста, зарегестрируйтесь/войдите, чтобы добавить товар в корзину",
    );
    return;
  }

  /** @type {Array<{id: number, count: number}>} */
  let cart;

  let id_bouquet = Number(id);
  let dataFromStorage = localStorage.getItem("cart");
  if (dataFromStorage == null) {
    //localStorage пока пустой
    //создаем аустой массив
    cart = [];
    let product = {
      id_bouquet: id_bouquet,
      count: 1,
    };
    cart.push(product); //записали объект в массив
  } else {
    // козина не была пустой
    // читаем данные которые уже были записаны
    cart = JSON.parse(dataFromStorage);
    let find_item = cart.find((item) => item.id_bouquet == id_bouquet);
    if (find_item != undefined) {
      //проверяем есть ли этот букет уже в корзине
      find_item.count += 1; // увеличиваем количество
    } else {
      // этого букета еще не было в корзине
      let product = {
        id_bouquet: id_bouquet,
        count: 1,
      };
      cart.push(product);
    }
  }

  //обновляем cart в localStorage
  localStorage.setItem("cart", JSON.stringify(cart));
  showNotification("Товар в корзине!");
}

export function getCartData() {
  // создадим Map для хранения данных корзины
  // Map - это коллекция ключ-значение, где ключи могут быть любого типа данных
  // таблица короче id-колв0
  let cartMap = new Map();
  let cart = JSON.parse(localStorage.getItem("cart")) || [];

  // заполняем Map данными из корзины
  cart.forEach((item) => {
    cartMap.set(item.id_bouquet, item.count);
  });

  return cartMap;
}

export const orderStatuses = {};

export async function fetchStatuses(){
  try{
    const response = await fetch("/order_statuses");
    if (!response.ok){
      throw new Error("Ошибка сервера: " + response.status);
    }
    const statuses = await response.json();
    statuses.forEach((status) => {
      orderStatuses[status.id_status] = status.name;
    });
  }catch(err){
    showNotification("Не удалось загрузить статусы заказов.");
    throw err;
  }
}

export function removeLoadings() {
  document.querySelectorAll(".loading").forEach((loading) => loading.remove());
}

export function clearCart() {
  localStorage.removeItem("cart");
}