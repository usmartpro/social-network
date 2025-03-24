-- +goose Up
WITH RECURSIVE
    names_arr AS (
        SELECT ARRAY[
                   'Александр', 'Максим', 'Артём', 'Михаил', 'Иван', 'Дмитрий', 'Даниил', 'Кирилл',
                   'Андрей', 'Егор', 'Никита', 'Илья', 'Алексей', 'Матвей', 'Тимофей', 'Роман',
                   'Владимир', 'Ярослав', 'Фёдор', 'Сергей', 'Глеб', 'Константин', 'Лев', 'Николай',
                   'Степан', 'Владислав', 'Павел', 'Георгий', 'Арсений', 'Денис', 'Виктор', 'Евгений',
                   'Марк', 'Олег', 'Юрий', 'Антон', 'Богдан', 'Василий', 'Захар', 'Семён', 'Пётр',
                   'Григорий', 'Станислав', 'Тимур', 'Артур', 'Давид', 'Игорь', 'Руслан', 'Вячеслав', 'Леонид'
                   ] AS names
    ),
    surnames_arr AS (
        SELECT ARRAY[
                   'Иванов', 'Смирнов', 'Кузнецов', 'Попов', 'Васильев', 'Петров', 'Соколов', 'Михайлов',
                   'Новиков', 'Фёдоров', 'Морозов', 'Волков', 'Алексеев', 'Лебедев', 'Семёнов', 'Егоров',
                   'Павлов', 'Козлов', 'Степанов', 'Николаев', 'Орлов', 'Андреев', 'Макаров', 'Никитин',
                   'Захаров', 'Зайцев', 'Соловьёв', 'Борисов', 'Яковлев', 'Григорьев', 'Романов', 'Воробьёв',
                   'Сергеев', 'Кузьмин', 'Фролов', 'Александров', 'Дмитриев', 'Королёв', 'Гусев', 'Киселёв',
                   'Ильин', 'Максимов', 'Поляков', 'Сорокин', 'Виноградов', 'Ковалёв', 'Белов', 'Медведев',
                   'Антонов', 'Тарасов'
                   ] AS surnames
    ),
    cities_arr AS (
        SELECT ARRAY[
                   'Москва', 'Санкт-Петербург', 'Новосибирск', 'Екатеринбург', 'Казань', 'Нижний Новгород',
                   'Челябинск', 'Самара', 'Омск', 'Ростов-на-Дону', 'Уфа', 'Красноярск', 'Пермь', 'Воронеж',
                   'Волгоград', 'Краснодар', 'Саратов', 'Тюмень', 'Тольятти', 'Ижевск'
                   ] AS cities
    ),
    hobbies_arr AS (
        SELECT ARRAY[
                   'Танцы', 'Пение', 'Шахматы', 'Бег', 'Йога', 'Фото', 'Рисунок', 'Вязание', 'Готовка', 'Рыбалка',
                   'Охота', 'Блог', 'Плавание', 'Велосипед', 'Туризм', 'Садоводство', 'Вышивка', 'Чтение', 'Игры', 'Анимация'
                   ] AS hobbies
    ),
    random_combinations AS (
        SELECT
            names[1 + floor(random() * array_length(names, 1))] AS name,
            surnames[1 + floor(random() * array_length(surnames, 1))] AS surname,
            DATE '2010-01-01' - (random() * (DATE '2010-01-01' - DATE '1930-01-01'))::int AS birth_date,
            cities[1 + floor(random() * array_length(cities, 1))] AS city,
            hobbies[1 + floor(random() * array_length(hobbies, 1))] AS hobby,
            1 AS n
        FROM names_arr, surnames_arr, cities_arr, hobbies_arr

        UNION ALL

        SELECT
            names[1 + floor(random() * array_length(names, 1))] AS name,
            surnames[1 + floor(random() * array_length(surnames, 1))] AS surname,
            DATE '2010-01-01' - (random() * (DATE '2010-01-01' - DATE '1930-01-01'))::int AS birth_date,
            cities[1 + floor(random() * array_length(cities, 1))] AS city,
            hobbies[1 + floor(random() * array_length(hobbies, 1))] AS hobby,
            n + 1
        FROM random_combinations, names_arr, surnames_arr, cities_arr, hobbies_arr
        WHERE n < 1000000
    )
INSERT INTO public.users (first_name, second_name, birthdate, city, biography, password)
SELECT
    name, surname, birth_date, city, hobby, md5(name || surname)
FROM random_combinations
LIMIT 500000;
