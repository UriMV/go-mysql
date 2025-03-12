// scripts.js
console.clear();

const loginBtn = document.getElementById('login');
const signupBtn = document.getElementById('signup');

loginBtn.addEventListener('click', (e) => {
    let parent = e.target.parentNode.parentNode;
    Array.from(e.target.parentNode.parentNode.classList).find((element) => {
        if (element !== "slide-up") {
            parent.classList.add('slide-up');
        } else {
            signupBtn.parentNode.classList.add('slide-up');
            parent.classList.remove('slide-up');
        }
    });
});

signupBtn.addEventListener('click', (e) => {
    let parent = e.target.parentNode;
    Array.from(e.target.parentNode.classList).find((element) => {
        if (element !== "slide-up") {
            parent.classList.add('slide-up');
        } else {
            loginBtn.parentNode.parentNode.classList.add('slide-up');
            parent.classList.remove('slide-up');
        }
    });
});

// scripts.js
document.addEventListener('DOMContentLoaded', function () {
    const navToggle = document.getElementById('nav-toggle');
    const navLinks = document.getElementById('nav-links');

    navToggle.addEventListener('click', function () {
        navLinks.classList.toggle('active');
    });
});