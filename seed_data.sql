TRUNCATE TABLE users RESTART IDENTITY;

INSERT INTO users (id, name, email, created_at, updated_at) VALUES
('6c37586f-86cb-410c-9d56-1afdc5bb6c5b', 'Hannah Jackson', 'hannah.jackson@example.com', '2022-03-07 22:24:27', '2022-06-17 22:24:27'),
('18f80bf3-9f52-4e98-8cde-2b1449195a5c', 'Fiona Thomas', 'fiona.thomas@mailservice.net', '2023-12-28 12:15:58', '2024-07-22 12:15:58'),
('5b98e3d6-4abe-4ea0-b86e-f6ef3c02c715', 'Charlie White', 'charlie.white@mailservice.net', '2022-05-12 14:03:57', '2023-02-01 14:03:57'),
('066a32da-a17f-43b6-ae05-bf1e7a4eb9c6', 'Alice Brown', 'alice.brown@example.com', '2022-10-02 10:18:24', '2022-11-14 10:18:24'),
('4782a391-a514-4a97-a52e-dae98d7d3067', 'George Smith', 'george.smith@mailservice.net', '2020-06-24 01:42:44', '2021-01-17 01:42:44'),
('38011fb2-9816-49c6-95d8-c4254468a3a3', 'Alice Smith', 'alice.smith@mailservice.net', '2024-11-30 20:38:57', '2025-06-11 20:38:57'),
('92e6e4d0-3cd4-4cac-b84c-2e9e2aa0714a', 'Ian Thomas', 'ian.thomas@mydomain.co', '2023-09-29 23:14:45', '2023-10-18 23:14:45'),
('5782d585-dacc-428b-a345-9428d79d4f7f', 'Ian Jackson', 'ian.jackson@demo.org', '2022-08-26 18:05:31', '2023-06-09 18:05:31'),
('5c9441d3-e776-4972-8910-694c5e0d5e27', 'George Jackson', 'george.jackson@example.com', '2021-08-03 16:15:59', '2021-10-07 16:15:59'),
('28152b2a-49e1-499e-a3dc-772160797edd', 'George Taylor', 'george.taylor@mydomain.co', '2020-04-04 07:51:35', '2020-10-19 07:51:35'),
('0c3a49c0-1f6b-4bca-9baa-d12ef05b4c5c', 'Julia White', 'julia.white@testmail.com', '2023-07-28 13:39:55', '2023-09-10 13:39:55'),
('4870d406-a658-4e37-a15a-85738f47ba1d', 'Bob Johnson', 'bob.johnson@example.com', '2020-03-09 16:51:20', '2020-09-06 16:51:20'),
('30bfe284-e49a-4669-b032-eac83089f2be', 'Ian Smith', 'ian.smith@example.com', '2023-08-13 15:09:30', '2023-10-14 15:09:30'),
('acdfe823-1d67-4bca-8d51-e342847d5987', 'Bob Jackson', 'bob.jackson@example.com', '2020-12-04 01:43:23', '2021-03-14 01:43:23'),
('3ddd59bc-b6c9-4bcf-8f94-20dc2d126706', 'Julia Brown', 'julia.brown@example.com', '2021-04-28 16:52:25', '2021-08-23 16:52:25'),
('c748f9f5-5794-42ca-ae12-a110b8eb75aa', 'Charlie Brown', 'charlie.brown@example.com', '2020-04-23 03:35:26', '2020-09-22 03:35:26'),
('dd53a770-5bfe-42a2-83bd-f2812d417d8e', 'Diana White', 'diana.white@example.com', '2024-06-13 12:34:56', '2025-06-03 12:34:56'),
('669cb4c5-c846-47d0-9183-6ec35d8aec90', 'George Anderson', 'george.anderson@testmail.com', '2023-04-28 20:43:49', '2023-06-26 20:43:49'),
('5d7afbfa-cc7f-423e-b47c-481c6b647a1f', 'Charlie Martin', 'charlie.martin@mailservice.net', '2024-08-18 01:10:17', '2025-03-04 01:10:17'),
('3aa4cdd7-3f2e-4c57-a323-c4d1caf71176', 'Alice Martin', 'alice.martin@demo.org', '2021-07-17 23:37:41', '2022-01-17 23:37:41');
