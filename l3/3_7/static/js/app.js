document.addEventListener("DOMContentLoaded", function () {
  // DOM Elements
  const loginSection = document.getElementById("login-section");
  const inventorySection = document.getElementById("inventory-section");
  const loginBtn = document.getElementById("login-btn");
  const signupBtn = document.getElementById("signup-btn");
  const addItemBtn = document.getElementById("add-item-btn");
  const itemFormSection = document.getElementById("item-form-section");
  const itemForm = document.getElementById("item-form");
  const cancelFormBtn = document.getElementById("cancel-form-btn");
  const inventoryTbody = document.getElementById("inventory-tbody");
  const historySection = document.getElementById("history-section");
  const historyTbody = document.getElementById("history-tbody");
  const closeHistoryBtn = document.getElementById("close-history-btn");
  const loginStatus = document.getElementById("login-status");

  // State
  let currentRole = null;

  // Event Listeners
  loginBtn.addEventListener("click", login);
  signupBtn.addEventListener("click", signup);
  addItemBtn.addEventListener("click", showAddForm);
  itemForm.addEventListener("submit", saveItem);
  cancelFormBtn.addEventListener("click", hideForm);
  closeHistoryBtn.addEventListener("click", hideHistory);

  // Initialize the app
  loadItems();

  // Functions
  function login() {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    currentRole = document.getElementById("role-select").value;

    // In a real app, you would send credentials to the server
    // For this demo, we'll simulate login success
    if (email && password) {
      currentUser = { email, role: currentRole };
      loginSection.classList.add("d-none");
      inventorySection.classList.remove("d-none");
      updateUIBasedOnRole();
      showStatus("Login successful!", "success");
    } else {
      showStatus("Please enter email and password", "error");
    }
  }

  function signup() {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const role = document.getElementById("role-select").value;

    if (email && password && role) {
      // For demo purposes, just simulate success
      currentUser = { email, role };
      loginSection.classList.add("d-none");
      inventorySection.classList.remove("d-none");
      updateUIBasedOnRole();
      showStatus("Signup successful!", "success");
    } else {
      showStatus("Please fill all fields", "error");
    }
  }

  function updateUIBasedOnRole() {
    const addItemBtn = document.getElementById("add-item-btn");
    const editButtons = document.querySelectorAll(".edit-btn");
    const deleteButtons = document.querySelectorAll(".delete-btn");
    const viewHistoryButtons = document.querySelectorAll(".view-history-btn");

    if (currentRole === "user") {
      addItemBtn.style.display = "none";
      editButtons.forEach((btn) => (btn.style.display = "none"));
      deleteButtons.forEach((btn) => (btn.style.display = "none"));
    } else if (currentRole === "manager") {
      addItemBtn.style.display = "inline-block";
      editButtons.forEach((btn) => (btn.style.display = "inline-block"));
      deleteButtons.forEach((btn) => (btn.style.display = "inline-block"));
      viewHistoryButtons.forEach((btn) => (btn.style.display = "inline-block"));
    } else if (currentRole === "admin") {
      addItemBtn.style.display = "inline-block";
      editButtons.forEach((btn) => (btn.style.display = "inline-block"));
      deleteButtons.forEach((btn) => (btn.style.display = "inline-block"));
      viewHistoryButtons.forEach((btn) => (btn.style.display = "inline-block"));
    }
  }

  async function loadItems() {
    try {
      const response = await fetch("/items");
      if (response.ok) {
        const items = await response.json();
        displayItems(items);
      } else {
        showStatus("Failed to load items", "error");
      }
    } catch (error) {
      showStatus("Error loading items: " + error.message, "error");
    }
  }

  function displayItems(items) {
    inventoryTbody.innerHTML = "";

    items.forEach((item) => {
      const row = document.createElement("tr");

      // Format price as currency
      const formattedPrice = item.price
        ? parseFloat(item.price).toFixed(2)
        : "0.00";

      row.innerHTML = `
                <td>${item.ID || item.id}</td>
                <td>${item.Name || item.name || ""}</td>
                <td>${item.Description || item.description || ""}</td>
                <td>${item.Quantity || item.quantity || 0}</td>
                <td>$${formattedPrice}</td>
                <td>
                    <button class="btn btn-primary btn-sm edit-btn me-1" data-id="${item.ID || item.id}">Edit</button>
                    <button class="btn btn-danger btn-sm delete-btn me-1" data-id="${item.ID || item.id}">Delete</button>
                    <button class="btn btn-info btn-sm view-history-btn" data-id="${item.ID || item.id}">History</button>
                </td>
            `;

      inventoryTbody.appendChild(row);
    });

    // Add event listeners to the new buttons
    document.querySelectorAll(".edit-btn").forEach((btn) => {
      btn.addEventListener("click", () =>
        showEditForm(parseInt(btn.dataset.id)),
      );
    });

    document.querySelectorAll(".delete-btn").forEach((btn) => {
      btn.addEventListener("click", () => deleteItem(parseInt(btn.dataset.id)));
    });

    document.querySelectorAll(".view-history-btn").forEach((btn) => {
      btn.addEventListener("click", () =>
        viewHistory(parseInt(btn.dataset.id)),
      );
    });

    updateUIBasedOnRole();
  }

  function showAddForm() {
    document.getElementById("form-title").textContent = "Add New Item";
    document.getElementById("item-id").value = "";
    document.getElementById("item-name").value = "";
    document.getElementById("item-description").value = "";
    document.getElementById("item-quantity").value = "";
    document.getElementById("item-price").value = "";
    itemFormSection.classList.remove("d-none");
  }

  function showEditForm(id) {
    // In a real app, we would fetch the item details by ID
    // For this demo, we'll just show an empty form
    document.getElementById("form-title").textContent = "Edit Item";
    document.getElementById("item-id").value = id;

    // For simplicity, we'll show empty fields - in a real app we'd pre-fill with existing data
    document.getElementById("item-name").value = "";
    document.getElementById("item-description").value = "";
    document.getElementById("item-quantity").value = "";
    document.getElementById("item-price").value = "";

    itemFormSection.classList.remove("d-none");
  }

  async function saveItem(e) {
    e.preventDefault();

    const id = document.getElementById("item-id").value;
    const name = document.getElementById("item-name").value;
    const description = document.getElementById("item-description").value;
    const quantity = parseInt(document.getElementById("item-quantity").value);
    const price = parseFloat(document.getElementById("item-price").value);

    const itemData = {
      name,
      description,
      quantity,
      price,
    };

    try {
      let response;
      if (id) {
        // Update existing item
        response = await fetch(`/items/${id}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(itemData),
        });
      } else {
        // Create new item
        response = await fetch("/items", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(itemData),
        });
      }

      if (response.ok) {
        itemFormSection.classList.add("d-none");
        loadItems();
        showStatus("Item saved successfully!", "success");
      } else {
        const errorData = await response.json();
        showStatus(
          `Failed to save item: ${errorData.error || "Unknown error"}`,
          "error",
        );
      }
    } catch (error) {
      showStatus("Error saving item: " + error.message, "error");
    }
  }

  function hideForm() {
    itemFormSection.classList.add("d-none");
  }

  async function deleteItem(id) {
    if (!confirm("Are you sure you want to delete this item?")) {
      return;
    }

    try {
      const response = await fetch(`/items/${id}`, {
        method: "DELETE",
      });

      if (response.ok) {
        loadItems();
        showStatus("Item deleted successfully!", "success");
      } else {
        showStatus("Failed to delete item", "error");
      }
    } catch (error) {
      showStatus("Error deleting item: " + error.message, "error");
    }
  }

  async function viewHistory(id) {
    try {
      // Fetch actual history records from the API
      const response = await fetch(`/items/${id}/history`);
      if (response.ok) {
        const history = await response.json();

        // Show history section
        historySection.classList.remove("d-none");

        // Display the history records
        displayHistory(history);
      } else {
        showStatus("Failed to load history", "error");
      }
    } catch (error) {
      showStatus("Error loading history: " + error.message, "error");
    }
  }

  function displayHistory(history) {
    historyTbody.innerHTML = "";

    history.forEach((record) => {
      const row = document.createElement("tr");

      // Format old and new values for display
      let oldValuesDisplay = record.old_values ? record.old_values : "N/A";
      let newValuesDisplay = record.new_values ? record.new_values : "N/A";

      // Limit the length of the JSON strings for display
      if (oldValuesDisplay.length > 50) {
        oldValuesDisplay = oldValuesDisplay.substring(0, 50) + "...";
      }
      if (newValuesDisplay.length > 50) {
        newValuesDisplay = newValuesDisplay.substring(0, 50) + "...";
      }

      row.innerHTML = `
                <td>${record.id}</td>
                <td>${record.action}</td>
                <td class="history-details">${oldValuesDisplay}</td>
                <td class="history-details">${newValuesDisplay}</td>
                <td>${record.changed_by}</td>
                <td>${new Date(record.changed_at).toLocaleString()}</td>
            `;

      historyTbody.appendChild(row);
    });
  }

  function hideHistory() {
    historySection.classList.add("d-none");
  }

  function showStatus(message, type) {
    loginStatus.textContent = message;
    loginStatus.className = `status-message ${type}`;

    // Auto-hide success messages after 3 seconds
    if (type === "success") {
      setTimeout(() => {
        loginStatus.textContent = "";
        loginStatus.className = "";
      }, 3000);
    }
  }
});
