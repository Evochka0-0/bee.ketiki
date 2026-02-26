import { addToCart, DeleteProduct, Home, getCartData, removeLoadings, showNotification, Cart, offerToLogIn, bouquet_images } from "./utils.js";


const buy_button = document.getElementById("buy_button");
const delete_button = document.getElementById("delete_button");
const name_container = document.getElementById("name_product_container");

let currentId;

async function ReviewWriter(currentId){
  // проверяем доступ к написанию отзыва
  let params = new URLSearchParams(document.location.search);
  const id_bouquet = Number(params.get("id"));

  if (!id_bouquet) return;
  try {
    const response = await fetch("/reviews/access", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(id_bouquet),
    });
  
    if (response.ok){
      const review_writer = document.getElementById('review_writer');
      //показываем 
      review_writer.classList.remove("review_writer_hidden");
    }
  } catch(err){
    console.error("Ошибка проверки доступа", err);
  }
  /**@type {HTMLSelectElement} */
  const reviw_select = document.getElementById('reviw_select');
  /** @type {HTMLTextAreaElement | null} */
  const textarea = document.getElementById('textarea');
  const but_send = document.getElementById('but_send');
  let message;
  let grade;

  but_send.addEventListener('click', () => {
    message = textarea.value.trim();
    grade = Number(reviw_select.value);
    if(message.length === 0 ){
      showNotification("Пустое поле, напишите свои впечатления о товаре");
      return;
    }
    if(message == "Напишите свои впечатления о товаре"){
      showNotification("Пустое поле, напишите свои впечатления о товаре");
      return;
    }

    if(grade == 0){
      showNotification("Не забудьте поставить оценку");
      return;
    }
    //данные для отпарвки на сервер
    let data = {
      id_bouquet : parseInt(id_bouquet),
      message: message,
      grade: grade
    }
    // теперь отправляем данные на сервер POST review
    fetch("/reviews/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    })
    .then(response => {
      if (!response.ok){
        throw new Error("Не удалось отправить отзыв" + response.status);
      }
      else{
        showNotification("Спасибо за отзыв!");
        window.location.reload();
      }
    })
    .catch(err => {
      showNotification("Произошла ошибка, попробуйте еще раз:", err.message);
      console.log("Ошибка", err.message, err);
    })
  })

}

function loadReviews(){
  let params = new URLSearchParams(document.location.search);
  const id_bouquet = Number(params.get("id"));

  if (!id_bouquet) return;

  fetch(`/reviews/${id_bouquet}`)
  .then(response => {
    if(!response.ok){
      throw new Error("Ошибка: " + response.status);
    }
    else{
      console.log("Получили response");
      removeLoadings();
      return response.json();
    }
  })
  .then((reviews) => {
    console.log("Получили response (.then((reviews) => {)");
    reviews.forEach((element) => {
      // для каждого отзыва 
      console.log("Создаем карточку");
      const review_container = document.createElement('div');
      review_container.className = "review_container";

      const grade_container = document.createElement('div');
      grade_container.className = "grade_container";
      let grade = element.grade;
      const grade_img = document.createElement('img');
      grade_img.className = "grade_stars";
      if(grade == 5){
        grade_img.src = "/images/5.png";
      }
      else if(grade == 4){
        grade_img.src = "/images/4.png";
      }
      else if(grade == 3){
        grade_img.src = "/images/3.png";
      }
      else if(grade == 2){
        grade_img.src = "/images/2.png";
      }
      else{
        grade_img.src = "/images/1.png";
      }

      const name_user = document.createElement('h3');
      name_user.textContent = `${element.last_name} ${element.first_name}`;

      grade_container.appendChild(name_user);
      grade_container.appendChild(grade_img);

      const main_container = document.createElement('div');
      main_container.className = "review_main";

      const text = document.createElement('h3');
      text.className = "review_text";
      text.textContent = element.message;

      main_container.appendChild(text);

      const created_at = document.createElement('span');
      created_at.className = "created_at";
      const date = new Date(element.created_at);
      created_at.textContent = date.toLocaleDateString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      });

      review_container.appendChild(grade_container);
      review_container.appendChild(main_container);
      review_container.appendChild(created_at);

      const reviews_div = document.getElementById('reviews');
      reviews_div.appendChild(review_container);

    });
  })
  .catch(err => {
    showNotification("Произошла какая-то ошибка");
    console.log(err.message, err);
  })
}

function loadInfo() {
  let params = new URLSearchParams(document.location.search);
  const id = Number(params.get("id"));

  if (!id) return;

  buy_button.onclick = () => {
    addToCart(id);
    updateQuantityInfo(id);
  };

  delete_button.onclick = () => {
    DeleteProduct(id);
    updateQuantityInfo(id);
  };

  updateQuantityInfo(id);

  fetch(`/bouquets/${id}`)
    .then((response) => response.json())
    .then((product) => {
      name_container.textContent = product.name;
      const img = document.getElementById("bouquet_image");
      bouquet_images(product, img);
      img.alt = product.name;
      document.getElementById("price_span").textContent = product.price;
      document.getElementById("description_container").textContent =
        product.description;
    });
}


function updateQuantityInfo(id) {
  const cartData = getCartData(); // Должен возвращать Map (id -> count)
  let count = cartData.get(id) || 0;

  if (count > 0) {
    buy_button.textContent = `В корзине ${count}`;
    buy_button.className = "incart_button";
    delete_button.style.display = "block";
  } else {
    buy_button.textContent = `В корзину`;
    buy_button.className = "buy_buttons";
    delete_button.style.display = "none";
  }
}

async function GetIdClient() {
  try {
    let role = "unauthorized";
    const response = await fetch("/auth_access");
    if (response.ok){
      const data = await response.json();
      currentId = data.user_id;
      //добавляем роль
      role = data.role;
    }

    switch (role){
      case "user":
        await ReviewWriter(currentId);
        break;
      default:
        return;
    }

    
  } catch (error) {
    console.log(error.message);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  offerToLogIn();
  loadInfo();
  loadReviews();
  GetIdClient();
  Home();
  Cart();
});
