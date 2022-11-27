// Utility functions:

const EL = (sel, par) => (par || document).querySelector(sel);
const ELS = (sel, par) => (par || document).querySelectorAll(sel);
const ELNew = (tag, prop) => Object.assign(document.createElement(tag), prop);

// App: Form: Add education:

const EL_add = EL(".add-button");
const EL_get = EL(".get-button");
const EL_fields = EL(".education-fields");

const addField = (name) => {
  const EL_field = ELNew("div", {
    className: "education-input-field row mb-2",
    innerHTML: `<div class="col-md-5">
      <input value="${
        name || ""
      }" type="text" class="form-control degree-name" placeholder="Degree name"/>
    </div>`
  });
  EL_fields.append(EL_field);
};

const getFieldsVaules = () => {
  return [...ELS(".degree-name", EL_fields)].map((el) => el.value);
};

// Events:

EL_add.addEventListener("click", () => {
  // Add a new field to DOM
  addField();
});

EL_get.addEventListener("click", () => {
  const fieldsValues = getFieldsVaules(); // Get all fields values into array
  console.log(fieldsValues);
  // TODO: do something with it
});

// Example: prepopulate with existing fields:

const existingFields = ["JavaScript Course", "HTML course"];
existingFields.forEach(addField); // That simple!