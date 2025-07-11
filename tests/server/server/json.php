<?php
header('Content-Type: application/json');

if (isset($_GET['syntax']) && str_contains($_GET['syntax'], "'")) {
    $response = [
        'status' => 'success',
        'data' => [
            'message' => 'This is a dummy response',
            'items' => [
                ['id' => 1, 'name' => 'Item 1'],
                ['id' => 2, 'name' => 'Item 2'],
                ['id' => 3, 'name' => 'Item 3'],
            ],
        ],
        'diff_key' => ['diff_value' => true] // Dummy diff 
    ];
} else {
    $response = [
    'status' => 'success',
    'data' => [
        'message' => 'This is a dummy response',
        'items' => [
            ['id' => 1, 'name' => 'Item 1'],
            ['id' => 2, 'name' => 'Item 2'],
            ['id' => 3, 'name' => 'Item 3'],
        ],
    ],
];
}
echo json_encode($response);
?>