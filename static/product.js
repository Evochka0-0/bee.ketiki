import { addToCart, DeleteProduct, Home, getCartData, removeLoadings, showNotification, Cart, offerToLogIn, bouquet_images } from "./utils.js";


const buy_button = document.getElementById("buy_button");
const delete_button = document.getElementById("delete_button");
const name_container = document.getElementById("name_product_container");

let currentId;
/**
 * 
 * @param {number} 
 * @param {number} 
 * @returns {Promise<string>} 
 */
async function initiatePayment(order_id, amount) {
    const response = await fetch("/payments", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id_order: order_id, amount: amount })
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Ошибка инициирования оплаты: " + response.status);
    }

    const data = await response.json();
    if (data.status !== "pending") {
        throw new Error("Оплата уже обработана");
    }
    return data.payment_ref;
}

/**
 * 
 * @param {number} user_id 
 * @param {number} id_bouquet 
 * @param {number} price 
 * @param {string} deadline 
 */
async function createQuickOrder(user_id, id_bouquet, price, deadline) {
    const orderData = {
        client_id: user_id,
        total_cost: price,
        deadline: deadline,
        items: [{ id_bouquet: id_bouquet, count: 1 }]
    };

    try {
        showNotification("Оформляем заказ...");
        
        const orderResponse = await fetch("/orders", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(orderData)
        });

        if (!orderResponse.ok) {
            const errorData = await orderResponse.json();
            throw new Error(errorData.error || "Ошибка создания заказа");
        }

        const orderDataResponse = await orderResponse.json();
        
        const payment_ref = await initiatePayment(orderDataResponse.id_order, price);
        
        window.location.replace(`/payment.html?id_order=${orderDataResponse.id_order}&amount=${price}&payment_ref=${payment_ref}`);
        
    } catch (error) {
        showNotification("Ошибка при оформлении заказа: " + error.message);
        console.error("Ошибка createQuickOrder:", error);
    }
}

function askForDeadlineAndCreate(user_id, id_bouquet, price) {
    const modal = document.createElement('div');
    modal.id = 'quick-order-deadline-modal';
    modal.className = 'modal-overlay';
    
    modal.innerHTML = `
        <div class="modal-window">
            <button id="close-deadline-modal" class="close-btn">✕</button>
            <h2>Выберите дату самовывоза</h2>
            <div class="deadline-selector">
                <input type="datetime-local" id="quick-order-deadline" class="deadline-input">
                <button id="confirm-deadline-btn" class="buy_buttons">Подтвердить и оплатить</button>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    const deadlineInput = document.getElementById('quick-order-deadline');
    const now = new Date();
    now.setHours(now.getHours() + 1);
    const minDateTime = now.toISOString().slice(0, 16);
    deadlineInput.min = minDateTime;
    deadlineInput.value = minDateTime;
    
    document.getElementById('confirm-deadline-btn').addEventListener('click', () => {
        const deadline = deadlineInput.value;
        if (!deadline) {
            showNotification("Выберите дату и время");
            return;
        }
        
        document.body.removeChild(modal);
        createQuickOrder(user_id, id_bouquet, price, deadline);
    });
    
    document.getElementById('close-deadline-modal').addEventListener('click', () => {
        document.body.removeChild(modal);
    });
}

function showGuestForm(id_bouquet, price, productName) {
    const modal = document.createElement('div');
    modal.id = 'quick-order-guest-modal';
    modal.className = 'modal-overlay';
    
    modal.innerHTML = `
        <div class="modal-window">
            <button id="close-guest-modal" class="close-btn">✕</button>
            <h2>Оформление заказа</h2>
            <p class="product-info">Товар: <strong>${productName}</strong></p>
            <p class="product-info">Сумма: <strong>${price.toLocaleString("ru-RU")} ₽</strong></p>
            
            <form id="quick-guest-form" class="guest-form">
                <div class="form-row">
                    <div class="form-group">
                        <label for="guest-lastname">Фамилия</label>
                        <input type="text" id="guest-lastname" required>
                    </div>
                    <div class="form-group">
                        <label for="guest-firstname">Имя</label>
                        <input type="text" id="guest-firstname" required>
                    </div>
                </div>
                
                <div class="form-group full-width">
                    <label for="guest-phone">Телефон</label>
                    <input type="tel" id="guest-phone" placeholder="+7 000 000 00 00" required>
                </div>
                
                <div class="form-group full-width">
                    <label for="guest-email">Электронная почта</label>
                    <input type="email" id="guest-email" placeholder="example@mail.com" required>
                </div>
                
                <div class="form-group full-width">
                    <label for="guest-deadline">Дата самовывоза</label>
                    <input type="datetime-local" id="guest-deadline" required>
                </div>
                
                <button type="submit" class="buy_buttons">Зарегистрироваться и оплатить</button>
            </form>
        </div>
    `;
    
    document.body.appendChild(modal);
    
    const deadlineInput = document.getElementById('guest-deadline');
    const now = new Date();
    now.setHours(now.getHours() + 1);
    const minDateTime = now.toISOString().slice(0, 16);
    deadlineInput.min = minDateTime;
    deadlineInput.value = minDateTime;
    
    document.getElementById('quick-guest-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const clientData = {
            last_name: document.getElementById('guest-lastname').value,
            first_name: document.getElementById('guest-firstname').value,
            phone: document.getElementById('guest-phone').value,
            email: document.getElementById('guest-email').value,
            password: 'temp_' + Math.random().toString(36).slice(2, 10)
        };
        
        const deadline = document.getElementById('guest-deadline').value;
        
        const submitBtn = e.target.querySelector('button[type="submit"]');
        const originalText = submitBtn.textContent;
        submitBtn.disabled = true;
        submitBtn.textContent = "Регистрируем...";
        
        try {
            const regResponse = await fetch('/clients', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(clientData)
            });
            
            if (!regResponse.ok) {
                const errorData = await regResponse.json();
                throw new Error(errorData.error || 'Ошибка регистрации');
            }
            
            const regData = await regResponse.json();
            
            document.body.removeChild(modal);
            
            await createQuickOrder(regData.id_client, id_bouquet, price, deadline);
            
        } catch (error) {
            showNotification('Ошибка: ' + error.message);
            submitBtn.disabled = false;
            submitBtn.textContent = originalText;
            console.error('Ошибка регистрации:', error);
        }
    });
    
    document.getElementById('close-guest-modal').addEventListener('click', () => {
        document.body.removeChild(modal);
    });
}

async function handleBuyNow() {
    const params = new URLSearchParams(document.location.search);
    const id_bouquet = Number(params.get("id"));
    
    if (!id_bouquet) {
        showNotification("Ошибка: не указан товар");
        return;
    }

    try {
        const response = await fetch(`/bouquets/${id_bouquet}`);
        if (!response.ok) throw new Error("Товар не найден");
        
        const product = await response.json();
        
        const authResponse = await fetch("/auth_access");
        
        if (authResponse.ok) {
            const userData = await authResponse.json();
            askForDeadlineAndCreate(userData.user_id, id_bouquet, product.price);
        } else {
            showGuestForm(id_bouquet, product.price, product.name);
        }
        
    } catch (error) {
        showNotification("Ошибка загрузки товара");
        console.error(error);
    }
}

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
    
    let data = {
      id_bouquet : parseInt(id_bouquet),
      message: message,
      grade: grade
    }
    
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
  const cartData = getCartData();
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
  
  const buyNowBtn = document.getElementById("buy_now_button");
  if (buyNowBtn) {
    buyNowBtn.addEventListener("click", handleBuyNow);
  }
  
  const deadlineInput = document.getElementById('datetime-local_deadline');
  if (deadlineInput) {
    const now = new Date();
    now.setHours(now.getHours() + 1);
    const minDateTime = now.toISOString().slice(0, 16);
    deadlineInput.min = minDateTime;
    deadlineInput.value = minDateTime;
  }
});