INSERT INTO users (id, username, email, password, created_at, updated_at)
VALUES (
           '98dd4aa0-4112-4040-9a95-e4130cf136d0',
           'testuser',
           'test@example.com',
           '$2a$10$Hy/5ZFnrV5wrk2dYlrVwTubi2hYRpnbcSg/uQQhinLptQsfA/caze',
           '2025-04-01 19:55:16.221+00',
           '2025-04-01 19:55:16.246+00'
       );
INSERT INTO business_customers (id, name, email, password, created_at, updated_at)
VALUES (
           '56e81a04-0162-464d-a42d-a27ca9156d25',
           'testbusinesscustomer',
           'test@example.com',
           '$2a$10$Yze4hlVdTUqxjbpCVilZaeMDNwjb59k1Hfn5cShIqnBwRiI5fxaXG',
           '2025-04-01 19:22:39.946+00',
           '2025-04-01 19:22:39.961+00'
       );
INSERT INTO service_templates (id, business_customer_id, name, description, fields_format, created_at, updated_at)
VALUES
    ('d6d18329-71f0-4da4-8bf4-46e402e17fe4', '56e81a04-0162-464d-a42d-a27ca9156d25', 'Premium Cleaning Service', 'A detailed cleaning service including windows and carpets.', NULL, '2025-04-01T19:30:21.478Z', '2025-04-01T19:30:21.489Z'),
    ('2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', '56e81a04-0162-464d-a42d-a27ca9156d25', 'Basic Cleaning Service', 'A standard cleaning service for homes and offices.', NULL, '2025-04-01T20:50:21.478Z', '2025-04-01T20:50:21.489Z'),
    ('5c1b8a77-53ea-4d55-a0d2-81433c7b4ef3', '56e81a04-0162-464d-a42d-a27ca9156d25', 'Deep Cleaning Service', 'An intensive cleaning service covering all areas.', NULL, '2025-04-01T20:50:21.500Z', '2025-04-01T20:50:21.512Z');

INSERT INTO questions (id, service_template_id, script_id, type, description, number, created_at, updated_at)
VALUES
    ('8be60c30-0be6-4adb-b4a6-11eb1a70282e', 'd6d18329-71f0-4da4-8bf4-46e402e17fe4', NULL, 'string', 'Generated question 1', 0, '2025-04-01T20:40:20.027025Z', '2025-04-01T20:40:20.027035Z'),
    ('87b62530-116c-40c1-bf8a-8285c36ac177', 'd6d18329-71f0-4da4-8bf4-46e402e17fe4', NULL, 'int', 'Generated question 2', 1, '2025-04-01T20:40:20.027126Z', '2025-04-01T20:40:20.027130Z'),
    ('687f6a4a-2f78-4141-a56a-6036a117d816', 'd6d18329-71f0-4da4-8bf4-46e402e17fe4', NULL, 'string', 'Generated question 3', 2, '2025-04-01T20:40:20.027237Z', '2025-04-01T20:40:20.027248Z'),
    ('2d135c19-d060-4873-8b47-01e71f731643', 'd6d18329-71f0-4da4-8bf4-46e402e17fe4', NULL, 'int', 'Generated question 4', 3, '2025-04-01T20:40:20.027365Z', '2025-04-01T20:40:20.027369Z'),
    ('30f1565f-d33e-4c79-8466-01054baf2232', 'd6d18329-71f0-4da4-8bf4-46e402e17fe4', NULL, 'string', 'Generated question 5', 4, '2025-04-01T20:40:20.027594Z', '2025-04-01T20:40:20.027599Z'),
    ('7adb4dc0-91b4-4142-9d2d-75c8d4a6015f', '2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', NULL, 'int', 'Generated question 6', 0, '2025-04-01T20:40:20.027941Z', '2025-04-01T20:40:20.027947Z'),
    ('7a58e4dd-ff76-4877-ba7d-f9910a585022', '2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', NULL, 'string', 'Generated question 7', 1, '2025-04-01T20:40:20.027999Z', '2025-04-01T20:40:20.028000Z'),
    ('ff380bb0-f7c1-479d-8f6c-0deae5f7d403', '2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', NULL, 'int', 'Generated question 8', 2, '2025-04-01T20:40:20.028017Z', '2025-04-01T20:40:20.028020Z'),
    ('9e57df38-3246-4869-a779-b52cb5e93c80', '2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', NULL, 'string', 'Generated question 9', 3, '2025-04-01T20:40:20.028033Z', '2025-04-01T20:40:20.028034Z'),
    ('72c62fe8-d984-4b72-bb79-e2fce35ed7ef', '2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', NULL, 'int', 'Generated question 10', 4, '2025-04-01T20:40:20.028043Z', '2025-04-01T20:40:20.028045Z'),
    ('7cf155e6-e6b7-4003-8ff1-c39fadc48359', '5c1b8a77-53ea-4d55-a0d2-81433c7b4ef3', NULL, 'string', 'Generated question 11', 0, '2025-04-01T20:40:20.028053Z', '2025-04-01T20:40:20.028055Z'),
    ('999d3792-c6b7-4fff-b48d-154663ebb105', '5c1b8a77-53ea-4d55-a0d2-81433c7b4ef3', NULL, 'int', 'Generated question 12', 1, '2025-04-01T20:40:20.028062Z', '2025-04-01T20:40:20.028064Z'),
    ('3efe4942-0ef5-4bda-bf19-2ae09646dcd5', '5c1b8a77-53ea-4d55-a0d2-81433c7b4ef3', NULL, 'string', 'Generated question 13', 2, '2025-04-01T20:40:20.028071Z', '2025-04-01T20:40:20.028081Z');
INSERT INTO filled_services (id, user_id, service_template_id, service_data, created_at, updated_at)
VALUES
    ('11111111-aaaa-bbbb-cccc-111111111111', '98dd4aa0-4112-4040-9a95-e4130cf136d0', 'd6d18329-71f0-4da4-8bf4-46e402e17fe4', NULL, '2025-04-02T12:00:00.000Z', '2025-04-02T12:00:00.000Z'),
    ('22222222-bbbb-cccc-dddd-222222222222', '98dd4aa0-4112-4040-9a95-e4130cf136d0', '2a8d9e66-b4f4-4c9a-9dd5-21e43f03c7cb', NULL, '2025-04-02T12:10:00.000Z', '2025-04-02T12:10:00.000Z'),
    ('33333333-cccc-dddd-eeee-333333333333', '98dd4aa0-4112-4040-9a95-e4130cf136d0', '5c1b8a77-53ea-4d55-a0d2-81433c7b4ef3', NULL, '2025-04-02T12:20:00.000Z', '2025-04-02T12:20:00.000Z');

INSERT INTO question_answered (id, question_id, filled_service_id, answer, number, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', '8be60c30-0be6-4adb-b4a6-11eb1a70282e', '11111111-aaaa-bbbb-cccc-111111111111', 'Tile', 0, '2025-04-02T10:00:00.000Z', '2025-04-02T10:00:00.000Z'),
    ('22222222-2222-2222-2222-222222222222', '87b62530-116c-40c1-bf8a-8285c36ac177', '11111111-aaaa-bbbb-cccc-111111111111', '1', 1, '2025-04-02T10:05:00.000Z', '2025-04-02T10:05:00.000Z'),
    ('33333333-3333-3333-3333-333333333333', '687f6a4a-2f78-4141-a56a-6036a117d816', '11111111-aaaa-bbbb-cccc-111111111111', 'Wood', 2, '2025-04-02T10:10:00.000Z', '2025-04-02T10:10:00.000Z'),
    ('44444444-4444-4444-4444-444444444444', '2d135c19-d060-4873-8b47-01e71f731643', '11111111-aaaa-bbbb-cccc-111111111111', '0', 3, '2025-04-02T10:15:00.000Z', '2025-04-02T10:15:00.000Z'),
    ('55555555-5555-5555-5555-555555555555', '30f1565f-d33e-4c79-8466-01054baf2232', '11111111-aaaa-bbbb-cccc-111111111111', 'Carpet', 4, '2025-04-02T10:20:00.000Z', '2025-04-02T10:20:00.000Z'),
    ('66666666-6666-6666-6666-666666666666', '7adb4dc0-91b4-4142-9d2d-75c8d4a6015f', '22222222-bbbb-cccc-dddd-222222222222', '1', 0, '2025-04-02T10:25:00.000Z', '2025-04-02T10:25:00.000Z'),
    ('77777777-7777-7777-7777-777777777777', '7a58e4dd-ff76-4877-ba7d-f9910a585022', '22222222-bbbb-cccc-dddd-222222222222', 'Vinyl', 1, '2025-04-02T10:30:00.000Z', '2025-04-02T10:30:00.000Z'),
    ('88888888-8888-8888-8888-888888888888', 'ff380bb0-f7c1-479d-8f6c-0deae5f7d403', '22222222-bbbb-cccc-dddd-222222222222', '0', 2, '2025-04-02T10:35:00.000Z', '2025-04-02T10:35:00.000Z'),
    ('99999999-9999-9999-9999-999999999999', '9e57df38-3246-4869-a779-b52cb5e93c80', '22222222-bbbb-cccc-dddd-222222222222', 'Laminate', 3, '2025-04-02T10:40:00.000Z', '2025-04-02T10:40:00.000Z'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '72c62fe8-d984-4b72-bb79-e2fce35ed7ef', '22222222-bbbb-cccc-dddd-222222222222', '1', 4, '2025-04-02T10:45:00.000Z', '2025-04-02T10:45:00.000Z'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '7cf155e6-e6b7-4003-8ff1-c39fadc48359', '33333333-cccc-dddd-eeee-333333333333', 'Concrete', 0, '2025-04-02T10:50:00.000Z', '2025-04-02T10:50:00.000Z'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '999d3792-c6b7-4fff-b48d-154663ebb105', '33333333-cccc-dddd-eeee-333333333333', '0', 1, '2025-04-02T10:55:00.000Z', '2025-04-02T10:55:00.000Z'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '3efe4942-0ef5-4bda-bf19-2ae09646dcd5', '33333333-cccc-dddd-eeee-333333333333', 'Marble', 2, '2025-04-02T11:00:00.000Z', '2025-04-02T11:00:00.000Z');

INSERT INTO scripts (id, name, script_code, created_at, updated_at)
VALUES
    ('a1b2c3d4-e5f6-7890-1234-56789abcdef0', 'Check More Letters Than Digits',
     'local input = arg[1]
      local letters, digits = 0, 0
      for c in input:gmatch(".") do
          if c:match("%a") then letters = letters + 1
          elseif c:match("%d") then digits = digits + 1 end
      end
      print(letters > digits)',
     '2025-04-01T21:00:00.000Z', '2025-04-01T21:00:00.000Z');
INSERT INTO scripts (id, name, script_code, created_at, updated_at)
VALUES
    ('b2c3d4e5-f678-9012-3456-789abcdef012', 'Check Even Integer',
     'local input = tonumber(arg[1])
      print(input and input % 2 == 0)',
     '2025-04-01T21:10:00.000Z', '2025-04-01T21:10:00.000Z');
INSERT INTO scripts (id, name, script_code, created_at, updated_at)
VALUES
    ('c3d4e5f6-7890-1234-5678-9abcdef01234', 'Check Positive Number',
     'local input = tonumber(arg[1])
      print(input and input > 0)',
     '2025-04-01T21:20:00.000Z', '2025-04-01T21:20:00.000Z');
INSERT INTO scripts (id, name, script_code, created_at, updated_at)
VALUES
    ('d4e5f6a7-8901-2345-6789-abcdef012345', 'Calculate Median',
     'local nums = {}
      for line in io.lines() do
          table.insert(nums, tonumber(line))
      end
      table.sort(nums)
      local n = #nums
      if n % 2 == 1 then
          print(nums[math.ceil(n / 2)])
      else
          print((nums[n / 2] + nums[n / 2 + 1]) / 2)
      end',
     '2025-04-01T21:30:00.000Z', '2025-04-01T21:30:00.000Z');
