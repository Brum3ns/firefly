<?php
header('Content-Type: application/xml; charset=utf-8');
echo '<?xml version="1.0" encoding="UTF-8"?>';
echo '<response>';
echo '  <status>success</status>';
echo '  <message>Dummy response from PHP server</message>';
if (isset($_GET['test']) && str_contains($_GET['test'], "'")) {
    echo '<diff>some diff</diff>';
}
echo '</response>';

?>