WITH
inserted_carriers AS (
INSERT INTO carriers (name, created_at, updated_at)
VALUES
    ('Nebulix Logística', NOW(), NOW()),
    ('RotaFácil Transportes', NOW(), NOW()),
    ('Moventra Express', NOW(), NOW())
    RETURNING id, name
    ),

    insert_policies AS (
INSERT INTO carrier_policies (carrier_id, region, estimated_days, price_per_kg, created_at, updated_at)
SELECT id, 'Sul', 4, 5.90, NOW(), NOW() FROM inserted_carriers WHERE name = 'Nebulix Logística'
UNION ALL
SELECT id, 'Sudeste', 4, 5.90, NOW(), NOW() FROM inserted_carriers WHERE name = 'Nebulix Logística'
UNION ALL
SELECT id, 'Sul', 7, 4.35, NOW(), NOW() FROM inserted_carriers WHERE name = 'RotaFácil Transportes'
UNION ALL
SELECT id, 'Sudeste', 7, 4.35, NOW(), NOW() FROM inserted_carriers WHERE name = 'RotaFácil Transportes'
UNION ALL
SELECT id, 'Centro-Oeste', 9, 6.22, NOW(), NOW() FROM inserted_carriers WHERE name = 'RotaFácil Transportes'
UNION ALL
SELECT id, 'Nordeste', 13, 8.00, NOW(), NOW() FROM inserted_carriers WHERE name = 'RotaFácil Transportes'
UNION ALL
SELECT id, 'Centro-Oeste', 7, 7.30, NOW(), NOW() FROM inserted_carriers WHERE name = 'Moventra Express'
UNION ALL
SELECT id, 'Nordeste', 10, 9.50, NOW(), NOW() FROM inserted_carriers WHERE name = 'Moventra Express'
    RETURNING *
    )
SELECT 1;
