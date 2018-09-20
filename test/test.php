<?php
passthru('clear');
ini_set('precision', 17);
$ua[]        = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/7046A194A';
$ua[]        = 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36';      
$type        = !empty($argv[1]) ? $argv[1] : "serial";
$iterations  = !empty($argv[2]) ? $argv[2] : 1;
$connections = !empty($argv[3]) ? $argv[3] : 1;
$test        = !empty($argv[4]) ? $argv[4] : 0;
$api         = !empty($argv[5]) ? $argv[5] : 'lb';

$stats = array(
  'requests'  => 0,
  'response'  => 0
);

echo "\e[104m ------------------------------------------------------------------------------------------------  \n\e[0m";
echo "\e[104m | USAGE: test.php TYPE[serial|parallel] ITERATIONS[#] CONNECTIONS[#] TEST. Default serial 1 1 0 | \n\e[0m";
echo "\e[104m | TYPE: serial: opens n connections in series, one afther the other, nicer on the server        | \n\e[0m";
echo "\e[104m | TYPE: parallel: opens n connections in parallel, all at once, call this 'MURRA!'              | \n\e[0m";
echo "\e[104m | ITERATIONS: how many times the processes will run                                             | \n\e[0m";
echo "\e[104m | CONNECTIONS PER SERIES: how many connections will be opened                                   | \n\e[0m";
echo "\e[104m | TEST sets the requests for test flow [1 | 0]                                                  | \n\e[0m";
echo "\e[104m ------------------------------------------------------------------------------------------------  \n\e[0m";
echo "\e[100m                                                                                                   \n";
echo "\e[100m CREATING {$connections} ".strtoupper($type)." CONNECTIONS AND RUNNING {$iterations} TIMES - TOTAL CONNECTIONS (".($iterations*$connections)."): \n";
echo "\e[100m                                                                                                   \n";
echo "\e[41m -------------------------------------------------------------------------------------------------  \n\e[0m";

$start = microtime(true);
$url   = 'http://127.0.0.1:8080/mutant';

include("dna.php");

//$func   = ($type === "serial") ? "serialRequests" : "parallelRequests";
$total  = 0;
for($i=0; $i < $iterations; $i++)
{
  if($type === "serial")
  {
    $current_time = serialRequests($url, $adns, $connections, $ua, $test, $stats);
  }
  else
  {
    $current_time = parallelRequests($url, $adns, $connections, $ua, $test, $stats);
  }
  $total += $current_time["total_time"];
  $curr = $i+1;
  echo "\e[41m -> Iteration $curr of $iterations done in " . $current_time['total_time'] . " / avg: " . $current_time['total_time']/$connections . " seconds \n\e[0m";
}
echo "\e[41m ------------------------------------------------------------------------------------------------  \n\e[0m";
echo "\e[0m                                                                                                    \n\e[0m";
echo "\e[0m [ TOTAL CONNECTIONS: " . ($iterations*$connections) . " ] \n";
echo "\e[0m [ TOTAL REQUESTS:        " . $stats["requests"] . " ] \n";
echo "\e[0m [ TOTAL RESPONSE:        " . $stats["response"] . " ] \n";
echo "\e[1m [ REQUESTS TOOK      " . $total . " SECONDS ] \n\e[0m";
echo "\e[1m [ SCRIPT TOOK        " . (microtime(true) - $start) . " SECONDS ] \n\e[0m";
echo "\e[0m                                                                                                    \n\e[0m";

function parallelRequests($url, $adns, $connections, $ua, $test, &$stats)
{
  $timing = ["total_time" => 0.0, "namelookup_time" => 0.0, "connect_time" => 0.0, "pretransfer_time" => 0.0];
  $curly = array();
  $mh = curl_multi_init();
  for($i=0; $i < $connections; $i++)
  {
    $data = array(
      'dna' => $adns[rand(0, 1041)]
    );
    $fields = json_encode($data);

    $curly[$i] = curl_init();
    curl_setopt($curly[$i], CURLOPT_URL,            $url);
    curl_setopt($curly[$i], CURLOPT_POST,           1);
    curl_setopt($curly[$i], CURLOPT_POSTFIELDS,     $fields);
    curl_setopt($curly[$i], CURLOPT_HTTPHEADER,     array('Content-Type: application/json', 'Content-Length: ' . strlen($fields)));
    curl_setopt($curly[$i], CURLOPT_FRESH_CONNECT,  1);
    curl_setopt($curly[$i], CURLOPT_HEADER,         0);
    curl_setopt($curly[$i], CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($curly[$i], CURLOPT_CONNECTTIMEOUT, 10);
    curl_setopt($curly[$i], CURLOPT_TIMEOUT,        10);
    curl_setopt($curly[$i], CURLOPT_USERAGENT,      $ua[rand(0,1)]);
    curl_multi_add_handle($mh, $curly[$i]);
  }
 
  $running = null;
  do 
  {
    curl_multi_exec($mh, $running);
  } 
  while($running > 0);
  
  foreach($curly as $id => $c)
  {
    $result = curl_multi_getcontent($c);
    // if($result !== false)
    // {
    //   $stats["bids"]++;
    //   notify($result, $ua, $test, $stats);
    // }
    $info = curl_getinfo($c);
    $timing["namelookup_time"] += (float)$info["namelookup_time"];
    $timing["pretransfer_time"] += (float)$info["pretransfer_time"] - (float)$info["namelookup_time"];
    $timing["connect_time"] += (float)$info["connect_time"] - (float)$info["starttransfer_time"];
    $timing["total_time"] += (float)$info["total_time"] - (float)$info["starttransfer_time"];
    curl_multi_remove_handle($mh, $c);
  }
  curl_multi_close($mh);
  $timing["namelookup_time"] /= $connections;
  $timing["pretransfer_time"] /= $connections;
  $timing["connect_time"] /= $connections;
  $timing["total_time"] /= $connections;
  
  return $timing;
}

function serialRequests($url, $adns, $connections, $ua, $test, &$stats)
{
  $timing = ["total_time" => 0.0, "namelookup_time" => 0.0, "connect_time" => 0.0, "pretransfer_time" => 0.0];
  for($i=0; $i < $connections; $i++)
  {
    $data = array(
      'dna' => $adns[rand(0, 1041)]
    );
    $fields = json_encode($data);
    $s = curl_init();
    curl_setopt($s, CURLOPT_URL,            $url);
    curl_setopt($s, CURLOPT_POST,           1);
    curl_setopt($s, CURLOPT_POSTFIELDS,     $fields);
    curl_setopt($s, CURLOPT_HTTPHEADER,     array('Content-Type: application/json', 'Content-Length: ' . strlen($fields)));
    curl_setopt($s, CURLOPT_HEADER,         0);
    curl_setopt($s, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($s, CURLOPT_CONNECTTIMEOUT, 10);
    curl_setopt($s, CURLOPT_TIMEOUT,        10);
    curl_setopt($s, CURLOPT_USERAGENT,      $ua[rand(0,1)]);
    $response = curl_exec($s);
    $timing = ["total_time" => 0.0, "namelookup_time" => 0.0, "connect_time" => 0.0, "pretransfer_time" => 0.0];
    $info = curl_getinfo($s);
    $timing["namelookup_time"] += (float)$info["namelookup_time"];
    $timing["pretransfer_time"] += (float)$info["pretransfer_time"] - (float)$info["namelookup_time"];
    $timing["connect_time"] += (float)$info["connect_time"] - (float)$info["starttransfer_time"];
    $timing["total_time"] += (float)$info["total_time"] - (float)$info["starttransfer_time"];
    curl_close($s);
    $stats["requests"]++;
    if($response !== false)
    {
      notify($response, $ua, $test, $stats);
    }
  }
  return $timing;
}

function notify($response, $ua, $test, &$stats)
{
  if(!empty($response))
  {
    //we need to parse this and fire a win notice and an impression (if possible)
    $data = json_decode($response);
    $stats['response']++;
  }
}
?>