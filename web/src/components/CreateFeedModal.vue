<script setup>
import { ref, defineEmits } from 'vue'

const emit = defineEmits(['close', 'create'])

const title = ref('')
const description = ref('')

function handleSubmit() {
  if (!title.value.trim()) return
  emit('create', title.value.trim(), description.value.trim())
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50" @click.self="emit('close')">
    <div class="card w-full max-w-md p-6">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Create New Feed</h2>

      <form @submit.prevent="handleSubmit">
        <div class="mb-4">
          <label for="title" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Title *
          </label>
          <input
            id="title"
            v-model="title"
            type="text"
            class="input"
            placeholder="My Podcast Feed"
            required
            autofocus
          />
        </div>

        <div class="mb-6">
          <label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Description
          </label>
          <textarea
            id="description"
            v-model="description"
            class="input"
            rows="3"
            placeholder="Optional description for your feed"
          ></textarea>
        </div>

        <div class="flex justify-end gap-3">
          <button type="button" @click="emit('close')" class="btn btn-secondary">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" :disabled="!title.trim()">
            Create Feed
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
