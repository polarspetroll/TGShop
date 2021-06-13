
function edit() {
  var pname = document.getElementById("pname").innerHTML;
  var pprice = document.getElementById("pprice").innerHTML;
  var pstat = document.getElementById("pstat").innerHTML;
  document.getElementById("editbutton").disabled = true;
  var form = document.createElement("form");
  form.setAttribute('method',"POST");
  form.setAttribute('id',"editform");
  var name = document.createElement("input");
  name.setAttribute('type',"text");
  name.setAttribute('name',"name");
  name.setAttribute('placeholder',"name");
  name.setAttribute('value', pname);
  var price = document.createElement("input");
  price.setAttribute('type', "text");
  price.setAttribute('name', "price");
  price.setAttribute('placeholder', "price");
  price.setAttribute('value', pprice);
  var status = document.createElement("label");
  status.innerHTML = "Availability";
  var select = document.createElement("select");
  select.setAttribute('name', "status");
  if (pstat == "in stock") {
    select.innerHTML = '<option value="1">In stock</option><option value="0">Out of stock</option>';
  } else {
    select.innerHTML = '<option value="0">Out of stock</option><option value="1">In stock</option>';
  }
  var submit = document.createElement("button");
  submit.innerHTML = "Save"
  var cancel = document.createElement("button");
  cancel.innerHTML = "Cancel";
  cancel.setAttribute('onclick', "Delete()");
  name.setAttribute('required', '');
  price.setAttribute('required', '');
  form.appendChild(name);
  form.appendChild(price);
  form.appendChild(status);
  form.appendChild(select);
  form.appendChild(submit);
  form.appendChild(cancel);
  document.getElementsByTagName('body')[0].appendChild(form);
}

function Delete() {
  form = document.getElementById("editform");
  form.parentNode.removeChild(form);
  document.getElementById("editbutton").disabled = false;
}

function Del() {
  if (!confirm('Are you sure you would like to remove this product?')) {
    throw '';
  }
  window.location.href='/del' + window.location.search;
}

function logout() {
  window.location.href = "/logout"
}
