window.addEventListener('DOMContentLoaded', () => {
  document.querySelector('#submit').addEventListener('click', () => {
    axios.post('/api', {
      name: document.querySelector('#name').value,
      email: document.querySelector('#email').value,
      content: document.querySelector('#content').value,
    }, {
      headers: {'Content-Type': 'application/json'},
      responseType: 'json',
    })
    .then((response) => {
      document.querySelector('#message').textContent = '';
      document.querySelector('#name').value = '';
      document.querySelector('#email').value = '';
      document.querySelector('#content').value = '';
    })
    .catch((error) => {
      document.querySelector('#message').textContent = error.response.data.error;
    });
  });
});
