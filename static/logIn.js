import { showNotification } from "./utils.js";

const link_p = document.getElementById("reg_link");
link_p.addEventListener('click', function() {
    window.location.replace('registration.html');
})


function loadInfo(){
    const input_phone = document.getElementById("phone").value;
    const input_password = document.getElementById("password").value;

    const inputData = {
        phone: input_phone,
        password: input_password
    }

    fetch('/login', {
        method: 'POST', // 2. Метод отправки
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(inputData)
        
    })
    .then(response =>{
        console.log("ответ получен");
        return response.json();
    })
    .then(data => {
        if(data.error){
            showNotification(data.error);
        }else{
            window.location.replace('index.html');
        }
    })
    .catch(err => {
        console.error("Сетевая ошибка:" + err);
        showNotification("Произошла ошибка при входе");
    });

}

const form = document.getElementById("Login_Form"); 
form.addEventListener('submit', function(event) {
    event.preventDefault(); 
    
    loadInfo();
});