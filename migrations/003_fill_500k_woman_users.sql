-- +goose Up
WITH RECURSIVE
    female_names_arr AS (
        SELECT ARRAY[
                   'Анна', 'Мария', 'Елена', 'Ольга', 'Наталья', 'Ирина', 'Светлана', 'Татьяна',
                   'Екатерина', 'Анастасия', 'Юлия', 'Александра', 'Дарья', 'Марина', 'Евгения',
                   'Виктория', 'Ксения', 'Людмила', 'Галина', 'Валентина', 'Лариса', 'Надежда',
                   'Вера', 'Алина', 'Маргарита', 'Валерия', 'Полина', 'София', 'Ангелина', 'Диана',
                   'Карина', 'Инна', 'Кристина', 'Яна', 'Алла', 'Лидия', 'Любовь', 'Зинаида',
                   'Вероника', 'Оксана', 'Тамара', 'Регина', 'Ульяна', 'Алёна', 'Элина', 'Эльвира',
                   'Василиса', 'Милана', 'Лилия', 'Снежана'
                   ] AS names
    ),
    female_surnames_arr AS (
        SELECT ARRAY[
                   'Иванова', 'Смирнова', 'Кузнецова', 'Попова', 'Васильева', 'Петрова', 'Соколова', 'Михайлова',
                   'Новикова', 'Фёдорова', 'Морозова', 'Волкова', 'Алексеева', 'Лебедева', 'Семёнова', 'Егорова',
                   'Павлова', 'Козлова', 'Степанова', 'Николаева', 'Орлова', 'Андреева', 'Макарова', 'Никитина',
                   'Захарова', 'Зайцева', 'Соловьёва', 'Борисова', 'Яковлева', 'Григорьева', 'Романова', 'Воробьёва',
                   'Сергеева', 'Кузьмина', 'Фролова', 'Александрова', 'Дмитриева', 'Королёва', 'Гусева', 'Киселёва',
                   'Ильина', 'Максимова', 'Полякова', 'Сорокина', 'Виноградова', 'Ковалёва', 'Белова', 'Медведева',
                   'Антонова', 'Тарасова'
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
        FROM female_names_arr, female_surnames_arr, cities_arr, hobbies_arr

        UNION ALL

        SELECT
            names[1 + floor(random() * array_length(names, 1))] AS name,
            surnames[1 + floor(random() * array_length(surnames, 1))] AS surname,
            DATE '2010-01-01' - (random() * (DATE '2010-01-01' - DATE '1930-01-01'))::int AS birth_date,
            cities[1 + floor(random() * array_length(cities, 1))] AS city,
            hobbies[1 + floor(random() * array_length(hobbies, 1))] AS hobby,
            n + 1
        FROM random_combinations, female_names_arr, female_surnames_arr, cities_arr, hobbies_arr
        WHERE n < 1000000
    )
INSERT INTO public.users (first_name, second_name, birthdate, city, biography, password)
SELECT
    name, surname, birth_date, city, hobby, md5(name || surname)
FROM random_combinations
LIMIT 500000;