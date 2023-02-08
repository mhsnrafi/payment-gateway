INSERT INTO cards (card_number, balance) VALUES
 ('1234567890123456', 100),
 ('2345678901234567', 200),
 ('3456789012345670', 300),
 ('4567890123456789', 400),
 ('5678901234567890', 500);


INSERT INTO frauds (card_number, amount, currency) VALUES
('1234567890123456', 0, 'USD'),
('2345678901234567', 0, 'EUR'),
('3456789012345670', 0, 'GBP'),
('4567890123456789', 125, 'CAD'),
('5678901234567890', 150, 'AUD');