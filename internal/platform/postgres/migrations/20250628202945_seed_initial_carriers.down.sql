DELETE FROM carrier_policies
WHERE carrier_id IN (
    SELECT id FROM carriers WHERE name IN ('Nebulix Logística', 'RotaFácil Transportes', 'Moventra Express')
);

DELETE FROM carriers
WHERE name IN ('Nebulix Logística', 'RotaFácil Transportes', 'Moventra Express');