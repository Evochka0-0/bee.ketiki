-- phpMyAdmin SQL Dump
-- version 3.5.1
-- http://www.phpmyadmin.net
--
-- Хост: 127.0.0.1
-- Время создания: Мар 20 2026 г., 22:44
-- Версия сервера: 5.5.25
-- Версия PHP: 5.3.13

SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- База данных: `bee.кетики1`
--

-- --------------------------------------------------------

--
-- Структура таблицы `base_color_palette`
--

CREATE TABLE IF NOT EXISTS `base_color_palette` (
  `id_base_color` int(11) NOT NULL AUTO_INCREMENT,
  `hex` varchar(7) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id_base_color`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=10 ;

--
-- Дамп данных таблицы `base_color_palette`
--

INSERT INTO `base_color_palette` (`id_base_color`, `hex`, `name`) VALUES
(1, 'FF0000', 'не указан'),
(2, 'FF0000', 'красный'),
(3, '4CB522', 'зеленый'),
(4, '2D24DD', 'синий'),
(5, 'DD24A9', 'розовый'),
(6, 'F4DA34', 'желтый'),
(7, 'FF9500', 'оранжевый'),
(8, '00C8FF', 'голубой'),
(9, 'A100FF', 'фиолетовый');

-- --------------------------------------------------------

--
-- Структура таблицы `bouquets`
--

CREATE TABLE IF NOT EXISTS `bouquets` (
  `id_bouquet` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  `description` text,
  `price` decimal(10,2) NOT NULL,
  `image_url` varchar(500) DEFAULT NULL,
  `reserve_image_url` varchar(500) NOT NULL COMMENT 'дополнительная картинка для витринных букетов',
  `id_base_color` int(11) NOT NULL,
  `type` varchar(50) NOT NULL DEFAULT 'usual',
  `id_occasion` int(11) NOT NULL COMMENT 'назначение',
  `id_structure` int(11) NOT NULL,
  PRIMARY KEY (`id_bouquet`),
  KEY `id_base_color` (`id_base_color`),
  KEY `occasion` (`id_occasion`),
  KEY `price` (`price`,`id_base_color`,`id_occasion`),
  KEY `id_occasion` (`id_occasion`),
  KEY `id_occasion_2` (`id_occasion`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=14 ;

--
-- Дамп данных таблицы `bouquets`
--

INSERT INTO `bouquets` (`id_bouquet`, `name`, `description`, `price`, `image_url`, `reserve_image_url`, `id_base_color`, `type`, `id_occasion`, `id_structure`) VALUES
(1, 'Нежные пионы', 'Пышный букет из нежных пионов, обернутые золотистой упаковочной бумагой', '2500.00', '/images/piones.png', '', 5, 'usual', 0, 0),
(2, 'Летнее небо', 'Авторский букет из нежных ромашек и ароматной лаванды.', '5100.00', '/images/romashki_lavanda.png', '', 6, 'usual', 0, 0),
(3, 'Альстромерии', 'Яркий и стойкий букет из 9 розовых альстромерий.', '3500.00', '/images/alstromerii.png', '', 5, 'usual', 0, 0),
(4, 'Белоснежный каприз', 'Пышный букет из крупных белых хризантем с кремово-белой лентой.', '2600.00', '/images/white.png', '', 6, 'usual', 1, 0),
(5, 'Розочки', 'Кустовые розы в классической крафтовой упаковке с красной лентой.', '3500.00', '/images/roses_craft.png', '', 5, 'usual', 0, 0),
(6, 'Орхидеи с жемчугом', 'Пышный букет из ярких фуксийно-розовых орхидей в подарочной упаковке.', '4000.00', '/images/fuksia_orhideia.png', '', 5, 'usual', 0, 0),
(7, 'Тюльпаны и эустома', 'Утонченный и элегантный минималистичный букет из нежных тюльпанов.', '2000.00', '/images/tulpani.png', '', 5, 'usual', 0, 0),
(8, 'Изюминка', 'Авторский, фактурный букет, который относится к современному стилю.', '5500.00', '/images/extra.png', '', 3, 'usual', 0, 0),
(9, 'Гжель с жемчугом', 'Фактурный букет с васильками, хлопковыми коробочками и жемчугом.', '3500.00', '/images/Vasilki.png', '', 8, 'usual', 0, 0),
(10, 'Хрустальное изящество', 'Букет из нежно-розовых калл, обернутый в белую матовую бумагу.', '2500.00', '/images/hrustal.png', '', 5, 'usual', 1, 0),
(11, 'Французские розы', 'Букет из изысканных французских садовых роз, обернутый в стильный крафт.', '6500.00', '/images/frahc_roses.png', '', 5, 'usual', 0, 0),
(12, 'Для зайчиков', 'Композиция из роз и лилий, украшенная зефиром и бусинами', '8500.00', '/images/carrots.png', '/images/carrots-0.png', 7, 'special', 0, 0),
(13, 'Ко дню влюбленных', 'Букет из хризантем и свежей клубники для вашей второй половинки', '8000.00', '/images/strowberryes.png', '/images/strowberryes-0.png', 2, 'special', 2, 0);

-- --------------------------------------------------------

--
-- Структура таблицы `bouquet_structure`
--

CREATE TABLE IF NOT EXISTS `bouquet_structure` (
  `id_bouquet_structure` int(11) NOT NULL AUTO_INCREMENT,
  `id_bouquet` int(11) NOT NULL,
  `id_flower` int(11) NOT NULL,
  PRIMARY KEY (`id_bouquet_structure`),
  KEY `id_bouquet` (`id_bouquet`),
  KEY `id_flower` (`id_flower`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=20 ;

--
-- Дамп данных таблицы `bouquet_structure`
--

INSERT INTO `bouquet_structure` (`id_bouquet_structure`, `id_bouquet`, `id_flower`) VALUES
(1, 1, 3),
(2, 2, 2),
(3, 2, 13),
(4, 3, 3),
(5, 4, 5),
(6, 5, 1),
(7, 6, 6),
(8, 7, 7),
(9, 7, 9),
(10, 8, 1),
(11, 8, 14),
(12, 9, 8),
(13, 9, 10),
(14, 10, 11),
(15, 11, 1),
(16, 12, 1),
(17, 12, 12),
(18, 12, 1),
(19, 13, 5);

-- --------------------------------------------------------

--
-- Структура таблицы `clients`
--

CREATE TABLE IF NOT EXISTS `clients` (
  `id_client` int(11) NOT NULL AUTO_INCREMENT,
  `last_name` varchar(100) NOT NULL,
  `first_name` varchar(100) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `email` varchar(191) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` varchar(20) NOT NULL DEFAULT 'user',
  PRIMARY KEY (`id_client`),
  UNIQUE KEY `phone` (`phone`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=24 ;

--
-- Дамп данных таблицы `clients`
--

INSERT INTO `clients` (`id_client`, `last_name`, `first_name`, `phone`, `email`, `password`, `role`) VALUES
(1, 'Сборщик', 'Ева', '+71111111111', 'zhuzhaev353@gmail.com', '$2a$10$HlYDw0dABz1a5IEcu6kY0O.iC3V1JCou6lArQG0xshXUljxUpXOsK', 'admin'),
(5, 'Цветочек', 'Димоооооооон', '+78005553535', 'privetik@gmail.com', '$2a$10$j2f5dGlsoESnDfjnadPHres6gQUoaqE9NKWacDw1b5fb.Pr3Que5W', 'user'),
(11, 'Иванов', 'Иван', '+70001212233', 'testIvan@gmai.com', '$2a$10$frTlBn.g6ZJWm0lT7RqOf.w5Zszzm4zPLI3.uz3jwCjAVGagbLONO', 'user'),
(14, 'Иванов', 'Иван', '+71112923344', 'testIv7an@gmai.com', '$2a$10$Wx9DkRvPLq/o6zeHbGKUA.kx5K32vqLYqWVHvUyJFhbF0iFytItGe', 'user'),
(15, 'Иванов', 'Иван', '+71112223349', 'test6_user@gmail.com', '$2a$10$DgOvwvs.BwyR2YSlJzXH8uCB1hTHDn8PKo2yp6Dh48shFhatMLQMa', 'user'),
(16, 'Шепелев', 'Егор', '89063441271', '616radical@gmail.com', '$2a$10$hEOUqS7Nq88AEptSlpVnfuvYTtbtZFQuEeRdURtYoCtVJq7eFEUGy', 'user'),
(17, 'Шепелева', 'Евка', '+71112221122', '6199adical@gmail.com', '$2a$10$kl/VQ3RhFXZSclK/lBRoDONVvewAM1bDBPkvMHTEBnx/9goV6ciUO', 'user'),
(19, 'Eva', 'Степа', '+71112223340', 'fgfgfgg@gmail.com', '$2a$10$Bybf3ql.hHDjURNl3SBCyu5e2rrZT2NXholr4Jph7mMSsnE83PQMC', 'user'),
(20, 'Гостев', 'Гость', '+79608208773', 'gostev@gmail.com', '$2a$10$sUEmqVkI7J5CR3xzxE0H0uAjHWpnP8JhCZTqWRYkidysfOoFYRIGG', 'user'),
(21, 'Eva', 'Димоооооооон', '+78005553539', 'zhuzh432353@gmail.com', '$2a$10$UziF8pzqpEVAcz2s6q6gKOT83tPX6Vh9.N6QZ4UlEXgHepKNidALS', 'user'),
(22, 'Ekkkkka', 'Димоооооооон', '+71456781122', 'zhu2222v353@gmail.com', '$2a$10$Cus.syp58EmhWq.pddNuIes.OUiAQNQI44Ghqb7SaUiBuWQGOffm.', 'user'),
(23, 'Шелепаева', 'Еварр', '+79994445511', 'zhuzh5555ev353@gmail.com', '$2a$10$IvHvovUP0erQoqL9/eOppOryZ8dXzN2txK8kMPKvDvPwysQx0Mk/q', 'user');

-- --------------------------------------------------------

--
-- Структура таблицы `flowers`
--

CREATE TABLE IF NOT EXISTS `flowers` (
  `id_flower` int(11) NOT NULL AUTO_INCREMENT,
  `name_flower` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id_flower`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=15 ;

--
-- Дамп данных таблицы `flowers`
--

INSERT INTO `flowers` (`id_flower`, `name_flower`) VALUES
(1, 'Розы'),
(2, 'Ромашки'),
(3, 'Пионы'),
(4, 'Альстромерии'),
(5, 'Хризантемы'),
(6, 'Орхидеи'),
(7, 'Тюльпаны'),
(8, 'Хлопок'),
(9, 'Эустома'),
(10, 'Васильки'),
(11, 'Каллы'),
(12, 'Лилии'),
(13, 'Лаванда'),
(14, 'Гвоздики');

-- --------------------------------------------------------

--
-- Структура таблицы `occasion`
--

CREATE TABLE IF NOT EXISTS `occasion` (
  `id_occassion` int(11) NOT NULL AUTO_INCREMENT,
  `ocassion_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id_occassion`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=4 ;

--
-- Дамп данных таблицы `occasion`
--

INSERT INTO `occasion` (`id_occassion`, `ocassion_name`) VALUES
(0, 'без назначения'),
(1, 'Свадебный'),
(2, 'День всех влюбленных');

-- --------------------------------------------------------

--
-- Структура таблицы `orderitems`
--

CREATE TABLE IF NOT EXISTS `orderitems` (
  `id_items` int(11) NOT NULL AUTO_INCREMENT,
  `id_order` int(11) NOT NULL,
  `id_bouquet` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  PRIMARY KEY (`id_items`),
  KEY `fk_items_order` (`id_order`),
  KEY `fk_items_bouquet` (`id_bouquet`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=56 ;

--
-- Дамп данных таблицы `orderitems`
--

INSERT INTO `orderitems` (`id_items`, `id_order`, `id_bouquet`, `quantity`) VALUES
(11, 8, 4, 1),
(12, 9, 6, 1),
(13, 9, 5, 1),
(14, 9, 2, 1),
(15, 9, 3, 1),
(16, 9, 1, 1),
(17, 10, 7, 1),
(18, 10, 8, 1),
(19, 10, 9, 1),
(20, 10, 10, 1),
(21, 10, 11, 1),
(22, 11, 11, 1),
(23, 12, 1, 1),
(24, 13, 3, 1),
(25, 14, 4, 1),
(26, 15, 11, 1),
(27, 16, 3, 1),
(28, 16, 4, 1),
(29, 17, 11, 1),
(30, 18, 9, 1),
(31, 19, 8, 1),
(32, 20, 4, 1),
(33, 21, 3, 1),
(34, 22, 2, 1),
(35, 23, 2, 1),
(36, 23, 3, 1),
(37, 23, 4, 1),
(38, 23, 6, 1),
(39, 23, 7, 1),
(40, 23, 7, 1),
(41, 23, 7, 1),
(42, 23, 5, 1),
(43, 24, 6, 3),
(44, 25, 3, 1),
(45, 25, 5, 2),
(46, 26, 13, 2),
(47, 26, 4, 1),
(48, 26, 12, 1),
(49, 27, 2, 1),
(50, 28, 2, 1),
(51, 29, 13, 1),
(52, 30, 13, 1),
(53, 31, 7, 1),
(54, 32, 7, 1),
(55, 33, 2, 1);

-- --------------------------------------------------------

--
-- Структура таблицы `orders`
--

CREATE TABLE IF NOT EXISTS `orders` (
  `id_order` int(11) NOT NULL AUTO_INCREMENT,
  `id_client` int(11) NOT NULL,
  `id_status` int(11) NOT NULL,
  `total_cost` decimal(10,2) NOT NULL,
  `payment_status` varchar(20) NOT NULL DEFAULT 'pending',
  `payment_ref` varchar(100) DEFAULT NULL,
  `deadline` datetime NOT NULL COMMENT 'дата когда клиент хочет получить товар',
  `pickup_datetime` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_order`),
  KEY `fk_orders_client` (`id_client`),
  KEY `fk_orders_status` (`id_status`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=34 ;

--
-- Дамп данных таблицы `orders`
--

INSERT INTO `orders` (`id_order`, `id_client`, `id_status`, `total_cost`, `payment_status`, `payment_ref`, `deadline`, `pickup_datetime`, `created_at`) VALUES
(8, 5, 4, '2600.00', 'paid', '6bbe82eaf57992aa82efd6bd', '0000-00-00 00:00:00', '2026-01-19 11:00:18', '2026-01-19 06:56:49'),
(9, 5, 4, '18600.00', 'paid', 'e198827ab7cd41d4eacea04c', '0000-00-00 00:00:00', '2026-01-19 11:00:21', '2026-01-19 06:59:01'),
(10, 5, 4, '20000.00', 'paid', '099e0381da07a6bdc3f6b462', '0000-00-00 00:00:00', '2026-01-19 11:00:23', '2026-01-19 06:59:25'),
(11, 5, 4, '6500.00', 'paid', '03a5f8ce78d92f4f9e04a9bc', '0000-00-00 00:00:00', '2026-02-03 19:03:55', '2026-01-19 07:09:55'),
(12, 5, 5, '2500.00', 'paid', 'bd68ee32f2d5d84c433ec772', '0000-00-00 00:00:00', NULL, '2026-01-19 07:13:47'),
(13, 5, 5, '3500.00', 'paid', '9fd50b09f64b9e9a9f2d2ba1', '0000-00-00 00:00:00', NULL, '2026-01-19 07:25:41'),
(14, 5, 5, '2600.00', 'paid', '962a655d07ed4da288fec66c', '0000-00-00 00:00:00', NULL, '2026-01-20 01:42:38'),
(15, 16, 5, '6500.00', 'paid', 'ae3df5a3c3e7846960251698', '0000-00-00 00:00:00', NULL, '2026-02-01 08:01:32'),
(16, 5, 3, '6100.00', 'paid', 'a236e7239762be0eafa65fa5', '2026-02-06 12:07:00', NULL, '2026-02-03 11:07:32'),
(17, 5, 3, '6500.00', 'paid', 'f41f270a0bc16a4bdf3fde3e', '2026-02-26 00:26:00', NULL, '2026-02-03 11:09:33'),
(18, 5, 3, '3500.00', 'paid', '9368d4e0fec58e10f527bb96', '2026-02-21 14:41:00', NULL, '2026-02-03 14:41:21'),
(19, 5, 4, '5500.00', 'paid', '4008233a38a633d2d0c1392e', '2026-02-04 16:00:00', '2026-02-26 21:04:37', '2026-02-04 07:00:24'),
(20, 5, 3, '2600.00', 'paid', '8ca5ee935057df6a83c07822', '2026-02-06 16:04:00', NULL, '2026-02-04 07:04:05'),
(21, 17, 3, '3500.00', 'paid', 'e9d5b830e17db17479ff1ad3', '2026-02-07 23:04:00', NULL, '2026-02-07 14:04:18'),
(22, 17, 3, '5100.00', 'paid', '9bba5944bf614448f181384d', '2026-02-07 23:11:00', NULL, '2026-02-07 14:11:34'),
(23, 5, 5, '20700.00', 'paid', '67613cc29caffdf3bfa02c72', '2026-02-12 16:51:00', NULL, '2026-02-11 16:51:58'),
(24, 5, 3, '12000.00', 'paid', '555f7ccec33a620f5df492e3', '2026-02-12 18:05:00', NULL, '2026-02-12 09:06:03'),
(25, 5, 4, '10500.00', 'paid', '37dceb5b4c4f24ec0d682a3d', '2026-02-15 19:39:00', '2026-02-26 21:04:30', '2026-02-12 10:39:45'),
(26, 5, 4, '27100.00', 'paid', 'fc1b4b93745bdca19e913a26', '2026-02-13 17:26:00', '2026-02-26 21:04:33', '2026-02-13 08:26:04'),
(27, 5, 1, '5100.00', 'paid', 'bec073b946525d8b4dd28351', '2026-02-17 12:01:00', NULL, '2026-02-17 03:01:10'),
(28, 5, 5, '5100.00', 'pending', 'fa47da5f2f373428ddee1810', '2026-02-25 13:13:00', NULL, '2026-02-25 04:13:49'),
(29, 5, 5, '8000.00', 'pending', '9d4a1359d3c5b8f9f72663a5', '2026-02-25 13:28:00', NULL, '2026-02-25 04:28:07'),
(30, 5, 1, '8000.00', 'paid', 'ff79c6bf32d5db62113f4e11', '2026-02-27 02:03:00', NULL, '2026-02-26 17:03:54'),
(31, 20, 1, '2000.00', 'paid', 'b8b2c5f38682d08373bd05b6', '2026-02-28 18:22:00', NULL, '2026-02-28 09:23:33'),
(32, 5, 1, '2000.00', 'pending', 'a378e95711133fb1450c0f0e', '2026-03-18 18:49:00', NULL, '2026-03-16 09:49:23'),
(33, 23, 1, '5100.00', 'paid', 'c7260d38bf9bb212a5e37785', '2026-03-16 20:56:00', NULL, '2026-03-16 15:57:27');

-- --------------------------------------------------------

--
-- Структура таблицы `orderstatuses`
--

CREATE TABLE IF NOT EXISTS `orderstatuses` (
  `id_status` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id_status`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=6 ;

--
-- Дамп данных таблицы `orderstatuses`
--

INSERT INTO `orderstatuses` (`id_status`, `name`) VALUES
(1, 'Новый'),
(2, 'Собирается'),
(3, 'Готов к выдаче'),
(4, 'Выдан'),
(5, 'Отменён');

-- --------------------------------------------------------

--
-- Структура таблицы `reviews`
--

CREATE TABLE IF NOT EXISTS `reviews` (
  `id_review` int(11) NOT NULL AUTO_INCREMENT,
  `id_client` int(11) NOT NULL,
  `id_bouquet` int(11) NOT NULL,
  `message` text NOT NULL,
  `grade` tinyint(4) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_review`),
  KEY `fk_reviews_client` (`id_client`),
  KEY `fk_reviews_bouquet` (`id_bouquet`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=8 ;

--
-- Дамп данных таблицы `reviews`
--

INSERT INTO `reviews` (`id_review`, `id_client`, `id_bouquet`, `message`, `grade`, `created_at`) VALUES
(1, 5, 1, 'урааа', 5, '2026-01-19 06:45:28'),
(2, 5, 3, 'ни то ни се', 3, '2026-01-19 06:45:51'),
(3, 5, 2, 'Ромашки вялые', 3, '2026-01-19 07:01:06'),
(4, 5, 4, 'Очень красивые, спасибо', 5, '2026-01-19 07:01:25'),
(5, 5, 5, 'Спасибо!', 4, '2026-01-19 07:01:48'),
(6, 5, 9, 'все хорошо,спасибо', 5, '2026-02-12 09:11:22'),
(7, 5, 7, 'красивый бантик папе понравилось', 1, '2026-02-25 04:40:38');

-- --------------------------------------------------------

--
-- Структура таблицы `sessions`
--

CREATE TABLE IF NOT EXISTS `sessions` (
  `id_session` int(11) NOT NULL AUTO_INCREMENT,
  `id_client` int(11) NOT NULL,
  `token` varchar(191) NOT NULL,
  `expires_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_session`),
  UNIQUE KEY `token` (`token`),
  KEY `fk_sessions_client` (`id_client`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=64 ;

--
-- Дамп данных таблицы `sessions`
--

INSERT INTO `sessions` (`id_session`, `id_client`, `token`, `expires_at`) VALUES
(50, 5, 'dfb689ac-0f39-45c3-818c-0098c31487cf', '2026-02-26 04:13:36'),
(51, 5, 'a629b72d-4283-4274-b71a-60d5d66a7de1', '2026-02-26 17:04:10'),
(52, 1, 'ee158d0e-8b4d-48fe-8954-e13c68b1f3d3', '2026-02-26 17:04:57'),
(53, 19, 'b8df5887-c6f6-4785-ac3e-9da8f1d0bb05', '2026-02-26 17:05:47'),
(54, 5, 'edc61fb0-449e-4164-8950-b917b6a23580', '2026-02-28 08:15:47'),
(55, 20, '971f3c44-9bb7-4cea-b087-7aba3224e893', '2026-02-28 09:24:15'),
(56, 21, '32934589-80e6-4b28-b83f-6cb3dda14d53', '2026-02-28 09:26:06'),
(57, 22, '4362a2ba-d78c-45a4-a69e-3389b240753b', '2026-02-28 09:26:57'),
(58, 5, 'bab6fd92-6ac8-4f17-9185-62876ee27011', '2026-03-16 09:01:50'),
(59, 5, '13842f98-dc1c-4b27-bd66-036f2757810d', '2026-03-16 09:02:18'),
(60, 5, '58221658-4448-45c4-9b07-6443f35df054', '2026-03-16 09:56:22'),
(61, 23, '371e748a-4a3e-439a-8b5b-fa2553296f58', '2026-03-16 16:12:44'),
(62, 5, 'c3bc88fd-21e6-4771-881f-f57e8a562c43', '2026-03-16 16:23:53'),
(63, 5, '7a6b80aa-e16d-4d9c-921e-4690c0cac76c', '2026-03-18 07:57:29');

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `bouquets`
--
ALTER TABLE `bouquets`
  ADD CONSTRAINT `bouquets_ibfk_2` FOREIGN KEY (`id_occasion`) REFERENCES `occasion` (`id_occassion`),
  ADD CONSTRAINT `bouquets_ibfk_1` FOREIGN KEY (`id_base_color`) REFERENCES `base_color_palette` (`id_base_color`);

--
-- Ограничения внешнего ключа таблицы `bouquet_structure`
--
ALTER TABLE `bouquet_structure`
  ADD CONSTRAINT `bouquet_structure_ibfk_1` FOREIGN KEY (`id_bouquet`) REFERENCES `bouquets` (`id_bouquet`),
  ADD CONSTRAINT `bouquet_structure_ibfk_2` FOREIGN KEY (`id_flower`) REFERENCES `flowers` (`id_flower`);

--
-- Ограничения внешнего ключа таблицы `orderitems`
--
ALTER TABLE `orderitems`
  ADD CONSTRAINT `fk_items_bouquet` FOREIGN KEY (`id_bouquet`) REFERENCES `bouquets` (`id_bouquet`),
  ADD CONSTRAINT `fk_items_order` FOREIGN KEY (`id_order`) REFERENCES `orders` (`id_order`) ON DELETE CASCADE;

--
-- Ограничения внешнего ключа таблицы `orders`
--
ALTER TABLE `orders`
  ADD CONSTRAINT `fk_orders_client` FOREIGN KEY (`id_client`) REFERENCES `clients` (`id_client`),
  ADD CONSTRAINT `fk_orders_status` FOREIGN KEY (`id_status`) REFERENCES `orderstatuses` (`id_status`);

--
-- Ограничения внешнего ключа таблицы `reviews`
--
ALTER TABLE `reviews`
  ADD CONSTRAINT `fk_reviews_bouquet` FOREIGN KEY (`id_bouquet`) REFERENCES `bouquets` (`id_bouquet`),
  ADD CONSTRAINT `fk_reviews_client` FOREIGN KEY (`id_client`) REFERENCES `clients` (`id_client`);

--
-- Ограничения внешнего ключа таблицы `sessions`
--
ALTER TABLE `sessions`
  ADD CONSTRAINT `fk_sessions_client` FOREIGN KEY (`id_client`) REFERENCES `clients` (`id_client`) ON DELETE CASCADE;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
