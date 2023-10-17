const messages = document.querySelector('#chat_room').childNodes;
console.log(messages);

messages.forEach((i) => {
  if (i.firstChild.firstChild.textContent === name) {
    i.classList.remove('bg-[#f0f0f1]', 'text-gray-900');
    i.classList.add('bg-indigo-500', 'text-white', 'self-end');
  }
});
