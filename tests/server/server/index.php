<?php
$title = "Firefly testserver";
$desc = "dynamic black box testing";

function randomStr($length) {
    $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $charactersLength = strlen($characters);
    $randomString = '';
    for ($i = 0; $i < $length; $i++) {
        $randomString .= $characters[random_int(0, $charactersLength - 1)];
    }
    return $randomString;
}

// CRLF
if ( isset($_GET['crlf']) ) {
    $crlf =  $_GET['crlf'];

    if ( str_contains($crlf, '"') === false ) {
        header('X-Deleted: true');
        header('X-Appear: APPEAR');
    } else {
        header('X-Normal: normalResponse');    
        header('X-Appear: false');    
    }
    
    header('Content-Type: text/plain');
    header("X-Custom: $crlf");
    header("Cache-Control: no-cache, must-revalidate");
}

function syntax() {
    if ( isset($_GET['syntax']) && str_contains($_GET['syntax'], "'")) {
        return "syntax error";
    }
}


function reflect() {
    if ( isset($_GET['reflect']) ) {
        return $_GET['reflect'];
    }
    return "reflect parameter is missing";
}


function reflectDynamic() {
    if ( isset($_GET['reflect']) && random_int(0, 2) == 1 ) {
        return htmlentities($_GET['reflect']);
    }
    return "reflect parameter is missing";
}


function randomness() {
    /*
    CSRF ~ 16-32 bytes
    SESSIONS ~ 32 bytes
    */

    if ( isset($_GET['randomness']) ) {
        return randomStr(16);
    }
    return "";
}

function disappear() {}
    if ( isset($_GET['disappear']) ) {
        $str = "disappear result: ";
        if ( str_contains($_GET['disappear'], '"') === false) {
            $str .= "APPEAR";
        } 
        echo "<p>$str</p>";
    }
?><?= reflect() ?>
<html>
<head>
<title><?= $title ?></title>
</head>
<body>
<h1><?= $title ?></h1>
<h3><?= $desc ?></h3>


<!-- behaviors -->
<div class="syntax">
    <?= syntax() ?>
</div>

<div class="crlf">
</div>

<!-- behaviors -->
<div class="transformation">
</div>

<div class="disappear">
    <?= disappear() ?>
</div>

<div class="reflect">
    <?= reflect() ?>
</div>

<div class="reflectdynamic">
    <?= reflectDynamic()?>
</div>

<div class="randomness">
    <?= randomness() ?>
</div>

<body>
</html><?= reflect() ?>