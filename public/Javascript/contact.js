// Get DOM

function  dataContact() {
    let name = document.getElementById("input-name").value;
    let email = document.getElementById("input-email").value;
    let phone = document.getElementById("input-phone").value;
    let subject = document.getElementById("input-selected").value;
    let message = document.getElementById("input-description").value;

    // return alert(name);
    // return alert(email);
    // return alert(phone);
    // return alert(select);
    // return alert(description);

    // Condition if the value is empty
    if(name == "") {
        return alert("Fill Your Name First");
    } else if(email == "") {
        return alert("Fill Your Email First");
    } else if(phone == "") {
        return alert("Fill Your Phone Number First");
    } else if(subject == "") {
        return alert("Choose Your Profession First");
    } else if(message == "") {
        return alert("Fill Your Message First");
    } 

     //mailto
     let a = document.createElement("a");
     a.href = `mailto:${email}?subject=${subject}&body=Hello, my name is ${name} | ${message}`;
     a.click()
}