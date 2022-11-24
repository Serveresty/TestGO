function start() {
  var betamount = document.getElementById("BetAmount");
  var range = document.getElementById("Range");
  var multiply = document.getElementById("Multiply");
  var winchance = document.getElementById("WinChance");
  var profit = document.getElementById("Profit");

  let data = {
    BetAmount: betamount.value,
    Range: range.value,
    Multiply: multiply.value,
    WinChance: winchance.value,
    Profit: profit.value,
    Number1: 0,
    Number2: 0,
    Number3: 0,
    Number4: 0,
    NotEnoughMoney: false,
    OutOfRange: false,
    OutOfMultiply: false,
    UnknownWinChance: false,
  };

  fetch("/dice", {
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    },
      method: "POST",
      body: JSON.stringify(data)
  }).then((response) => {
      response.text().then(function (data) {
          let result = JSON.parse(data);
          console.log(result)

          if (result['NotEnoughMoney'] == true) {
            document.getElementById("NoMoney").innerHTML = "No Money";
          } else {
            document.getElementById("NoMoney").innerHTML = "";
          }

          if (result['OutOfRange'] == true) {
            alert("Wrong Range");
          }

          if (result['OutOfMultiply'] == true) {
            alert("Wrong Multiply");
          }

          if (result['UnknownWinChance'] == true) {
            alert("Wrong Win Chance");
          }


          betamount.value = result['BetAmount'];
          range.value = result['Range'];
          multiply.value = result['Multiply'];
          winchance.value = result['WinChance'];
          profit.value = result['Profit'];
          num1 = result['Number1'];
          num2 = result['Number2'];
          num3 = result['Number3'];
          num4 = result['Number4'];

          roller = document.getElementById("resultDice");

          roller.textContent = num1.toString() + num2.toString() + num3.toString() + num4.toString();
      });
  }).catch((error) => {
      console.log(error)
  });
}

function RangeEdit(range) {
  var wrg_rng = document.getElementById("wrg-rng");

  if (range.toString().length > 5){
    range = parseFloat(range.toString().substring(0, 5));
  }
  if (range > 94 || range < 0.01) {
    wrg_rng.innerHTML = "Wrong Range";
    return;
  } else {
    wrg_rng.innerHTML = " ";
  }
  var winchance = document.getElementById("WinChance");
  winchance.value = range;
}

function MultiplyEdit(multy) {
  var wrg_mtp = document.getElementById("wrg-mtp");
  if (multy > 9500 || multy < 1.0106) {
    wrg_mtp.innerHTML = "Wrong Multiply";
    return;
  } else {
    wrg_mtp.innerHTML = "";
  }

  var bet = document.getElementById("BetAmount").value;
  document.getElementById("Profit").value = bet * multy - bet;
}

function WinChanceEdit(win) {
  var wrg_wch = document.getElementById("wrg-wch");

  if (win.toString().length > 5){
    win = parseFloat(win.toString().substring(0, 5));
  }
  if (win > 94 || win < 0.01) {
    wrg_wch.innerHTML = "Wrong Chance";
    return;
  } else {
    wrg_wch.innerHTML = " ";
  }
  var range = document.getElementById("Range");
  range.value = win;
}

function BetEdit(bet) {
  var multy = document.getElementById("Multiply").value;
  document.getElementById("Profit").value = bet * multy - bet;
}
