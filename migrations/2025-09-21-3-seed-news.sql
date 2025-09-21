-- Seed base users (passwords are demo-only)
INSERT INTO users (email, name, password)
VALUES
    ('m.ashinov@example.com', 'Михаил Аршинов', 'password'),
    ('vas.prog@example.com', 'Василий Программист', 'password'),
    ('li.son@example.com', 'Ли Сон', 'password'),
    ('irina@example.com', 'Ирина Иванова', 'password'),
    ('alex@example.com', 'Алексей Смирнов', 'password')
ON CONFLICT (email) DO NOTHING;

-- Seed categories used on the site
INSERT INTO categories (title, alias)
VALUES
    ('Еда', 'eda'),
    ('Животные', 'zhivotnye'),
    ('Машины', 'mashiny'),
    ('Спорт', 'sport'),
    ('Музыка', 'muzyka'),
    ('Технологии', 'tehnologii'),
    ('Прочее', 'prochee')
ON CONFLICT (alias) DO NOTHING;

-- Seed 20 news items
INSERT INTO news (title, alias, preview, content, user_id, createdat, views)
VALUES
    ('Открытие сезона байдарок', 'otkrytie-sezona-baidarok', 'Старт сезона на воде: маршруты и советы.', 'Полный текст: советы по выбору снаряжения и маршрутов для новичков и опытных.', (SELECT id FROM users WHERE email='m.ashinov@example.com'), NOW() - INTERVAL '30 days', 120),
    ('Выбор ноутбука для задач', 'vybor-noutbuka-dlya-zadach', 'На что смотреть при покупке ноутбука.', 'Полный текст: процессоры, память, накопители, экраны и нюансы автономности.', (SELECT id FROM users WHERE email='vas.prog@example.com'), NOW() - INTERVAL '27 days', 220),
    ('Создание автомобилей с автопилотом', 'sozdanie-avtomobiley-s-avtopilotom', 'Как работает автопилот и его будущее.', 'Полный текст: сенсоры, карты, ML-модели, безопасность и правовые аспекты.', (SELECT id FROM users WHERE email='li.son@example.com'), NOW() - INTERVAL '25 days', 340),
    ('Как быстро приготовить вкусный обед', 'kak-bystro-prigotovit-vkusnyy-obed', '15 минут и обед готов.', 'Полный текст: простые ингредиенты, баланс белков и углеводов, экономия времени.', (SELECT id FROM users WHERE email='irina@example.com'), NOW() - INTERVAL '22 days', 410),
    ('Уход за домашними животными', 'uhod-za-domashnimi-zhivotnymi', 'Базовые правила ухода.', 'Полный текст: питание, прогулки, ветеринарные осмотры и дрессировка.', (SELECT id FROM users WHERE email='alex@example.com'), NOW() - INTERVAL '20 days', 180),
    ('Музыкальные тренды 2025', 'muzykalnye-trendy-2025', 'Что слушают в этом году.', 'Полный текст: жанры, плейлисты, новые артисты и популярные платформы.', (SELECT id FROM users WHERE email='m.ashinov@example.com'), NOW() - INTERVAL '18 days', 520),
    ('Новые технологии в образовании', 'novye-tehnologii-v-obrazovanii', 'EdTech ускоряет обучение.', 'Полный текст: адаптивные системы, LLM, AR/VR и аналитика.', (SELECT id FROM users WHERE email='vas.prog@example.com'), NOW() - INTERVAL '17 days', 260),
    ('Как начать бегать', 'kak-nachat-begat', 'Первые шаги в беге без травм.', 'Полный текст: техника, дыхание, план тренировок на 4 недели.', (SELECT id FROM users WHERE email='li.son@example.com'), NOW() - INTERVAL '15 days', 300),
    ('Лучшие снеки для офиса', 'luchshie-sneki-dlya-ofisa', 'Полезные перекусы на работе.', 'Полный текст: орехи, батончики, фрукты, питьевой режим и кофеин.', (SELECT id FROM users WHERE email='irina@example.com'), NOW() - INTERVAL '14 days', 140),
    ('Путешествие на озеро', 'puteshestvie-na-ozero', 'Как подготовиться к поездке.', 'Полный текст: маршрут, снаряжение, безопасность и бюджет.', (SELECT id FROM users WHERE email='alex@example.com'), NOW() - INTERVAL '13 days', 200),
    ('Роботы в быту', 'roboty-v-bytu', 'Чем уже помогают роботы дома.', 'Полный текст: пылесосы, lawn-боты, помощники и голосовые интерфейсы.', (SELECT id FROM users WHERE email='m.ashinov@example.com'), NOW() - INTERVAL '12 days', 330),
    ('Безопасная езда зимой', 'bezopasnaya-ezda-zimoy', 'Шины, дистанция, манёвры.', 'Полный текст: тормозной путь, ABS, и поведение на льду.', (SELECT id FROM users WHERE email='vas.prog@example.com'), NOW() - INTERVAL '11 days', 275),
    ('Советы по фотографии', 'sovety-po-fotografii', 'Композиция и свет.', 'Полный текст: правило третей, работа с тенью, RAW и постобработка.', (SELECT id FROM users WHERE email='li.son@example.com'), NOW() - INTERVAL '10 days', 195),
    ('ТОП‑5 игр 2025', 'top-5-igr-2025', 'Игры, которые стоит попробовать.', 'Полный текст: геймплей, производительность, мультиплеер и сюжеты.', (SELECT id FROM users WHERE email='irina@example.com'), NOW() - INTERVAL '9 days', 610),
    ('Уход за собакой летом', 'uhod-za-sobakoy-letom', 'Жара и активность.', 'Полный текст: лапы, питьё, длинные прогулки и безопасность у воды.', (SELECT id FROM users WHERE email='alex@example.com'), NOW() - INTERVAL '8 days', 165),
    ('Электронная коммерция: тренды', 'elektronnaya-kommerciya-trendy', 'Что меняется в eCommerce.', 'Полный текст: маркетплейсы, логистика, платежи и персонализация.', (SELECT id FROM users WHERE email='m.ashinov@example.com'), NOW() - INTERVAL '7 days', 455),
    ('Автокросс в городе', 'avtokross-v-gorode', 'Соревнования и безопасность.', 'Полный текст: трассы, техника, страховка и требования.', (SELECT id FROM users WHERE email='vas.prog@example.com'), NOW() - INTERVAL '6 days', 280),
    ('Гастро‑фестиваль в городе', 'gastro-festival-v-gorode', 'Новые вкусы и блюда.', 'Полный текст: локальные фермеры, дегустации и мастер‑классы.', (SELECT id FROM users WHERE email='li.son@example.com'), NOW() - INTERVAL '5 days', 390),
    ('Клубная музыка ночью', 'klubnaya-muzyka-nochyu', 'Сеты и атмосфера.', 'Полный текст: диджеи, оборудование и культура клуба.', (SELECT id FROM users WHERE email='irina@example.com'), NOW() - INTERVAL '4 days', 510),
    ('Как создать подкаст', 'kak-sozdat-podkast', 'Идея, запись, публикация.', 'Полный текст: микрофоны, монтаж, хостинг и продвижение.', (SELECT id FROM users WHERE email='alex@example.com'), NOW() - INTERVAL '3 days', 440)
ON CONFLICT (alias) DO NOTHING;

-- Link news to categories
-- 1
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='otkrytie-sezona-baidarok'), (SELECT id FROM categories WHERE alias='sport')) ON CONFLICT DO NOTHING;
-- 2
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='vybor-noutbuka-dlya-zadach'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;
-- 3
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='sozdanie-avtomobiley-s-avtopilotom'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='sozdanie-avtomobiley-s-avtopilotom'), (SELECT id FROM categories WHERE alias='mashiny')) ON CONFLICT DO NOTHING;
-- 4
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='kak-bystro-prigotovit-vkusnyy-obed'), (SELECT id FROM categories WHERE alias='eda')) ON CONFLICT DO NOTHING;
-- 5
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='uhod-za-domashnimi-zhivotnymi'), (SELECT id FROM categories WHERE alias='zhivotnye')) ON CONFLICT DO NOTHING;
-- 6
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='muzykalnye-trendy-2025'), (SELECT id FROM categories WHERE alias='muzyka')) ON CONFLICT DO NOTHING;
-- 7
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='novye-tehnologii-v-obrazovanii'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;
-- 8
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='kak-nachat-begat'), (SELECT id FROM categories WHERE alias='sport')) ON CONFLICT DO NOTHING;
-- 9
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='luchshie-sneki-dlya-ofisa'), (SELECT id FROM categories WHERE alias='eda')) ON CONFLICT DO NOTHING;
-- 10
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='puteshestvie-na-ozero'), (SELECT id FROM categories WHERE alias='prochee')) ON CONFLICT DO NOTHING;
-- 11
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='roboty-v-bytu'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;
-- 12
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='bezopasnaya-ezda-zimoy'), (SELECT id FROM categories WHERE alias='mashiny')) ON CONFLICT DO NOTHING;
-- 13
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='sovety-po-fotografii'), (SELECT id FROM categories WHERE alias='prochee')) ON CONFLICT DO NOTHING;
-- 14
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='top-5-igr-2025'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;
-- 15
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='uhod-za-sobakoy-letom'), (SELECT id FROM categories WHERE alias='zhivotnye')) ON CONFLICT DO NOTHING;
-- 16
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='elektronnaya-kommerciya-trendy'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;
-- 17
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='avtokross-v-gorode'), (SELECT id FROM categories WHERE alias='sport')) ON CONFLICT DO NOTHING;
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='avtokross-v-gorode'), (SELECT id FROM categories WHERE alias='mashiny')) ON CONFLICT DO NOTHING;
-- 18
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='gastro-festival-v-gorode'), (SELECT id FROM categories WHERE alias='eda')) ON CONFLICT DO NOTHING;
-- 19
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='klubnaya-muzyka-nochyu'), (SELECT id FROM categories WHERE alias='muzyka')) ON CONFLICT DO NOTHING;
-- 20
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='kak-sozdat-podkast'), (SELECT id FROM categories WHERE alias='muzyka')) ON CONFLICT DO NOTHING;
INSERT INTO news_categories (news_id, category_id) VALUES ((SELECT id FROM news WHERE alias='kak-sozdat-podkast'), (SELECT id FROM categories WHERE alias='tehnologii')) ON CONFLICT DO NOTHING;


