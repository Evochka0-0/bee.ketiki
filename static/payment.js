import { Home, showNotification, clearCart } from "./utils.js";

/**
 * Подтверждает "оплату" на сервере.
 * @param {string} orderId - ID заказа.
 * @param {string} paymentRef - Уникальный идентификатор платежа.
 * @returns {Promise<object>} - Данные об успешной оплате.
 */
async function confirmPayment(orderId, paymentRef) {
    const response = await fetch("/payments/confirm", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            id_order: parseInt(orderId),
            payment_ref: paymentRef,
        }),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `Ошибка подтверждения: ${response.status}`);
    }

    return response.json();
}

async function handlePaymentPage() {
    // авторизован ли пользователь
    try {
        const authResponse = await fetch("/auth_access");
        if (!authResponse.ok) throw new Error("Не авторизован");
    } catch (error) {
        window.location.href = "registration.html";
        return;
    }

    // параметры из URL (order_id, amount, payment_ref)
    const params = new URLSearchParams(window.location.search);
    const orderId = params.get("id_order");
    const amount = params.get("amount");
    const paymentRef = params.get("payment_ref");

    if (!orderId || !amount || !paymentRef) {
        document.getElementById("payment_main").innerHTML =
            "<h1>Ошибка: Не найдены данные для проведения оплаты.</h1><p>Пожалуйста, вернитесь в корзину и попробуйте снова.</p>";
        return;
    }

    document.getElementById("order_id_display").textContent = orderId;
    document.getElementById("order_amount_display").textContent = `${parseFloat(
        amount,
    ).toLocaleString("ru-RU")} ₽`;

    const confirmButton = document.getElementById("confirm_payment_button");
    const loader = document.getElementById("payment_loader");
    const successMessage = document.getElementById("payment_success");
    const paymentContainer = document.querySelector(".payment_container");

    confirmButton.addEventListener("click", async () => {
        confirmButton.disabled = true;

        document.getElementById("payment_title").classList.add("hidden");
        document.querySelector(".order_summary").classList.add("hidden");
        document.querySelector(".fake_payment_form").classList.add("hidden");
        confirmButton.classList.add("hidden");
        loader.classList.remove("hidden");

        try {
            //Имитируем задержку
            await new Promise((resolve) => setTimeout(resolve, 1500));

            await confirmPayment(orderId, paymentRef);

            clearCart();

            loader.classList.add("hidden");
            successMessage.classList.remove("hidden");

            setTimeout(() => {
                window.location.replace('/home.html');
            }, 4000);
        } catch (error) {
            showNotification(`Ошибка оплаты: ${error.message}`);
            // Возвращаем всё как было
            loader.classList.add("hidden");
            document.getElementById("payment_title").classList.remove("hidden");
            document.querySelector(".order_summary").classList.remove("hidden");
            document.querySelector(".fake_payment_form").classList.remove("hidden");
            confirmButton.classList.remove("hidden");
            confirmButton.disabled = false;
        }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    handlePaymentPage();
    Home(); 
});
