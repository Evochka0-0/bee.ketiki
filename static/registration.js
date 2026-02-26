import { showNotification } from "./utils.js";

const link_p = document.getElementById("logIn_link");
link_p.addEventListener('click', function() {
    window.location.replace('logIn.html');
})

function loadInfo(){
    const last_name_input = document.getElementById("lastname").value;
    const first_name_input = document.getElementById("firstname").value;
    const phone_input = document.getElementById("phone").value;
    const email_input = document.getElementById("email").value;
    const password_input = document.getElementById("password").value; 

    const clientData = {
        last_name: last_name_input,
        first_name: first_name_input,
        phone: phone_input,
        email: email_input,
        password: password_input
    };


    fetch('/clients', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(clientData)
    })
    .then(async response => {
        if (response.ok){
            return response.text();
        }
        else {
            const errorData = await response.json();
            throw new Error(errorData.error);
        }
    })
    .then(responseText => {
        showNotification(responseText);
        setTimeout(() => {
            window.location.replace('home.html');
        }, 1000);
    })
    .catch(error => {
        let text = error.message;
        if (text.includes("Duplicate entry")) {
            if (text.includes("'phone'")) {
                showNotification("Пользователь с таким телефоном уже существует!");
                return;
            }
            if (text.includes("'email'")) {
                showNotification("Эта электронная почта уже занята!");
                return;
            }
        }
        showNotification("Ошибка: " + text);
        console.error("Технические детали:", error);
    })
};

const form = document.getElementById("Form");

form.addEventListener('submit', function(event) {
    event.preventDefault(); 
    loadInfo(); 
});
