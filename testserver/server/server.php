<?php
$title = "Firefly testserver";
$desc = "dynamic black box testing";

// CRLF
if ( isset($_GET['crlf']) ) {
    $crlf =  $_GET['crlf'];

    if ( str_contains($crlf, '"') === false ) {
        echo "HEADER => APPEAR";
        header('X-Deleted: true');
        header('X-Appear: APPEAR');
    } else {
        header('X-Appear: false');    
    }
    
    header('Content-Type: text/plain');
    header("X-Custom: $crlf");
    header("Cache-Control: no-cache, must-revalidate");
}

function xss() {}

function ssti() {}

function disappear() {}
    if ( isset($_GET['disappear']) ) {
        $str = "disappear result: ";
        if ( str_contains($_GET['disappear'], '"') === false) {
            $str .= "APPEAR";
        } 
        echo "<p>$str</p>";
    }
?>


<html>
<head>
<title><?= $title ?></title>
</head>
<body>
<h1><?= $title ?></h1>
<h3><?= $desc ?></h3>


<!-- vulnerabilities -->
<div class="xss">
</div>

<div class="ssti">
</div>

<div class="sqli">
</div>

<div class="crlf">
</div>

<!-- behaviors -->
<div class="transformation">
</div>

<!-- dynamic content -->
<div class="dynamic">

</div>

<div class="disappear">
    <?= disappear() ?>
</div>

<body>
</html>