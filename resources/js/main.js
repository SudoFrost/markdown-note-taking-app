import Editor from '@toast-ui/editor';
import { loadIcons } from './icons'
loadIcons()
import '@toast-ui/editor/dist/toastui-editor.css'; // Editor's Style
const editor = new Editor({
  el: document.querySelector('#editor'),
  height: '100%',
  initialEditType: 'markdown',
  previewStyle: 'vertical',
  placeholder: '# Write your note here',
});

const TEMPLATES = Array.from(document.getElementsByTagName('template'))
const NOTE_TEMPLATE = TEMPLATES.find(t => t.id === 'note-template')
const NOTES = document.getElementById('notes')
const FORMS = Array.from(document.getElementsByTagName("form"))
const NOTE_FORM = FORMS.find((e) => e.id == "note-form")
const CREATE_NOTE = document.getElementById("create-note")
CREATE_NOTE.addEventListener("click", () => {
  NOTE_FORM.reset()
  NOTE_FORM.dataset.action = '/api/notes'
  NOTE_FORM.dataset.method = 'POST'
  editor.setMarkdown('')
  editor.focus()
})
NOTE_FORM.addEventListener("submit", (e) => e.preventDefault())

editor.addHook('change', () => {
  NOTE_FORM.getElementsByTagName('textarea')[0].value = editor.getMarkdown()
})

document.getElementById('search-notes').addEventListener('click', LoadNotes)

async function GetNotes(categoryId) {
  let url = new URL(window.location)
  url.pathname = '/api/notes'
  if (categoryId) {
    url.searchParams.set('category', categoryId)
  }

  let search = document.getElementById('search').value
  if (search) {
    url.searchParams.set('q', search)
  }

  const res = await fetch(url)
  return await res.json()
}

async function GetNote(id) {
  const res = await fetch('/api/notes/' + id)
  return await res.json()
}

async function EditNote(id) {
  const note = await GetNote(id)

  const formElement = document.getElementById('note-form')
  formElement.querySelector('input[name="title"]').value = note.Title
  formElement.querySelector('textarea[name="content"]').value = note.Content
  formElement.querySelector('select[name="category"]').value = note.CategoryID || ''
  editor.setMarkdown(note.Content)

  NOTE_FORM.dataset.action = '/api/notes/' + id
  NOTE_FORM.dataset.method = 'PUT'
}

async function DeleteNote(id) {
  await fetch('/api/notes/' + id, { method: 'DELETE' })
  LoadNotes()
}

function RenderNotes(notes) {
  NOTES.innerHTML = ''
  notes.forEach(note => {
    const noteElement = NOTE_TEMPLATE.content.cloneNode(true)
    noteElement.querySelector('li').dataset.id = note.ID
    noteElement.querySelector('a').textContent = note.Title
    noteElement.getElementById('edit-note').onclick = () => EditNote(note.ID)
    noteElement.getElementById('delete-note').onclick = () => DeleteNote(note.ID)
    NOTES.appendChild(noteElement)
  })
  loadIcons()
}

async function LoadNotes() {
  const categoryId = NOTES.dataset.category
  RenderNotes(await GetNotes(categoryId))
}

document.getElementById("save-note").addEventListener("click", async function () {
  const method = NOTE_FORM.dataset.method
  const action = NOTE_FORM.dataset.action
  const formData = new FormData(NOTE_FORM)
  const res = await fetch(action, {
    method,
    body: formData,
  })
  if (res.ok) {
    LoadNotes()
    alert(method == "post" ? "Note Created" : "Note Saved")
  }
})

const CATEGORY_TEMPLATE = TEMPLATES.find(t => t.id === 'category-template')
const CATEGORY_FORM_TEMPLATE = TEMPLATES.find(t => t.id === 'category-form-template')

function RenderCategories(categories) {
  const categoriesElement = document.getElementById('categories')
  const ActiveCategory = NOTES.dataset.category
  const newElements = []
  const noCategoryElement = CATEGORY_TEMPLATE.content.cloneNode(true).querySelector('li')
  noCategoryElement.querySelector('a').textContent = 'No Category'
  noCategoryElement.querySelector('li > div').remove()
  noCategoryElement.dataset.id = ""
  newElements.push(noCategoryElement)
  categories.forEach(category => {
    const categoryElement = CATEGORY_TEMPLATE.content.cloneNode(true).querySelector('li')
    categoryElement.dataset.id = category.ID
    categoryElement.querySelector('a').textContent = category.Name
    categoryElement.querySelector('#edit').onclick = () => EditCategory(category.ID)
    categoryElement.querySelector('#delete').onclick = () => DeleteCategory(category.ID)
    newElements.push(categoryElement)
  })
  categoriesElement.innerHTML = ''
  newElements.forEach(category => {
    category.querySelector('a').onclick = () => {
      SetActiveCategory(category.dataset.id)
      return false
    }
    if (category.dataset.id == ActiveCategory) category.dataset.active = true
    categoriesElement.appendChild(category)
  })
  loadIcons()
}

async function GetCategories() {
  const res = await fetch('/api/categories')
  return await res.json()
}

function UpdateNoteFormCategories(categories) {
  const selectElement = NOTE_FORM.querySelector('select')
  const options = Array.from(selectElement.querySelectorAll('option[value]')).reduce((acc, option) => {
    if (option.value) acc[option.value] = option
    return acc
  }, {})
  selectElement.innerHTML = ''
  const updatedOptions = []
  categories.forEach(category => {
    let option = options[category.ID]
    if (!option) {
      option = document.createElement('option')
    }
    option.value = category.ID
    option.textContent = category.Name
    updatedOptions.push(option)
  })
  const defaultOption = document.createElement('option')
  defaultOption.value = ''
  defaultOption.textContent = 'No Category'
  selectElement.append(defaultOption, ...updatedOptions)
}

async function LoadCategories() {
  const categories = await GetCategories()
  RenderCategories(categories)
  UpdateNoteFormCategories(categories)
}

document.getElementById("create-category").onclick = function () {
  const categoriesElement = document.getElementById('categories')
  let createCategoryFormElement = categoriesElement.querySelector('li:has(form[data-action="/api/categories"])')
  if (!createCategoryFormElement) {
    createCategoryFormElement = CATEGORY_FORM_TEMPLATE.content.cloneNode(true).querySelector('li')
    categoriesElement.prepend(createCategoryFormElement)
  }
  createCategoryFormElement.querySelector('input[name="name"]').value = ""
  createCategoryFormElement.querySelector('form').dataset.action = '/api/categories'
  createCategoryFormElement.querySelector('form').dataset.method = 'POST'
  createCategoryFormElement.querySelector('#cancel').onclick = createCategoryFormElement.remove.bind(createCategoryFormElement)
  createCategoryFormElement.querySelector('#save').onclick = async function () {
    const res = await fetch('/api/categories', {
      method: 'POST',
      body: new FormData(createCategoryFormElement.querySelector('form')),
    })
    if (res.ok) {
      LoadCategories()
      alert("Category Created")
    }
  }
  loadIcons()
}

function EditCategory(id) {
  const categoriesElement = document.getElementById('categories')
  const originalCategoryElement = categoriesElement.querySelector(`li[data-id="${id}"]`)
  const editCategoryFormElement = CATEGORY_FORM_TEMPLATE.content.cloneNode(true).querySelector('li')

  editCategoryFormElement.dataset.id = id
  const form = editCategoryFormElement.querySelector('form')
  const inputName = form.querySelector('input[name="name"]')
  inputName.value = originalCategoryElement.querySelector('a').textContent
  form.dataset.action = `/api/categories/${id}`
  form.dataset.method = 'PUT'
  form.querySelector('#cancel').onclick = () => {
    editCategoryFormElement.replaceWith(originalCategoryElement)
    loadIcons()
  }
  form.querySelector('#save').onclick = async function () {
    const response = await fetch(form.dataset.action, {
      method: form.dataset.method,
      body: new FormData(form),
    })
    if (response.ok) {
      LoadCategories()
      alert("Category Saved")
    }
  }

  originalCategoryElement.replaceWith(editCategoryFormElement)
  loadIcons()
}

function SetActiveCategory(id) {
  const categories = document.querySelectorAll('#categories > li[data-id]')
  categories.forEach(category => {
    category.dataset.active = category.dataset.id == id
  })
  NOTES.dataset.category = id
  LoadNotes()
}

async function DeleteCategory(id) {
  await fetch('/api/categories/' + id, {
    method: 'DELETE'
  })
  LoadCategories()
}
LoadCategories()
LoadNotes()

