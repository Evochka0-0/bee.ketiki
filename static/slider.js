import { showNotification } from "./utils.js";

let currentIndex = 0;
let count_slides;

let slides_array = new Array();

export function arrowsListener(){
    const back = document.getElementById("back");
    const forward = document.getElementById("forward");

    forward.addEventListener('click', () => {
        currentIndex ++ ;
        moveSlide(currentIndex, 1);
    })
    
    back.addEventListener('click', () => {
        currentIndex -- ;
        moveSlide(currentIndex, 1);
    })
}

function moveSlide(index, transform){
    const track = document.getElementById("track_container");
    if (transform == 1){//простое листание с анимацией
        track.style.transition = "transform 0.7s ease";
    }else{
        track.style.transition = "none";
    }
    const offset = index * -100;
    track.style.transform = `translateX(${offset}%)`;
}

function CreateSlide(bouquet){
    const slide = document.createElement('div');
    slide.className = "spacial_bouq_slide";

    slide.style.backgroundImage = `url("${bouquet.image_url}")`;
    slide.style.minHeight = "84vh";
    slide.style.backgroundSize = "cover";
    slide.style.backgroundPosition = "center";

    const name = document.createElement('h1');
    name.className = "spacial_bouq_name";
    name.textContent = bouquet.name;
    console.log(bouquet.dominate_color);
    name.style.color = bouquet.dominate_color;

    const special_buy_button = document.createElement('button');
    special_buy_button.className = "special_buy_buttons";
    special_buy_button.textContent = `${bouquet.price} | В корзину`;

    const datas = document.createElement('div');
    datas.className = "name_button";
    datas.style.display = 'flex';
    datas.style.justifyContent = 'center';

    datas.appendChild(name);
    datas.appendChild(special_buy_button);

    slide.appendChild(datas);

    return slide;
}

export async function loadSlides(){
    const track = document.getElementById("track_container");

    //загружаем все букеты из коллекции
    try{
        const bouquets = await fetch("/bouquets?type=special").then((response) =>
            response.json(),
        );
        
        let bouquets_length = bouquets.length;
        let first_clone;
        let last_clone;
        bouquets.forEach(bouquet => {
            let slide = CreateSlide(bouquet);
            slides_array.push(slide);
            track.append(slide);
        });
        first_clone = slides_array[0].cloneNode(true);
        last_clone = slides_array[bouquets_length - 1].cloneNode(true);
        track.prepend(last_clone);
        track.append(first_clone);
        track.addEventListener('transitionend', () => {
            track.style.transition = 'none';
            if (currentIndex >= bouquets_length + 1){//если на последнем слайде(клон первого)
                currentIndex = 1;
                moveSlide(currentIndex, 0);
            }else if(currentIndex <= 0){
                currentIndex = bouquets_length;
                moveSlide(currentIndex, 0);
            }
        });
        currentIndex = 1;
        moveSlide(currentIndex, 0);
    }catch(err){
        console.error(err.message, err);
        showNotification("Не удалось загрузить особую коллекцию");
    }
}