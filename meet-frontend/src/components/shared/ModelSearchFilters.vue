<script setup lang="ts">
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import { Field, Form } from 'vee-validate'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import * as yup from 'yup'
import { FilterSearchModel } from '../../api/domain/Model'
import { GeneralUtils } from '../../utils/GeneralUtils'

const emits = defineEmits<{(e: 'filters-changed', filters: FilterSearchModel): void }>()

const { t } = useI18n({ useScope: 'global' })

const initialValues = {
  nickName: null,
  countryCode: null,
  cityName: null,
  zodiacSignCode: null
}

const schema = yup.object({
  nickName: yup.string().notRequired(),
  countryCode: yup.string().notRequired(),
  cityName: yup.string().notRequired(),
  zodiacSignCode: yup.string().notRequired()
})

const countries = ref([
  { name: 'Argentina', code: 'AR' },
  { name: 'Chile', code: 'CL' },
  { name: 'Perú', code: 'PE' },
  { name: 'Bolivia', code: 'BOL' },
  { name: 'Uruguay', code: 'UR' },
  { name: 'Colombia', code: 'COL' },
  { name: 'Ecuador', code: 'ECU' },
  { name: 'Venezuela', code: 'VE' }
])

const zodiacSigns = ref([
  { name: 'Aries', code: 'ARI' },
  { name: 'Tauro', code: 'TAU' },
  { name: 'Geminis', code: 'GEM' },
  { name: 'Cáncer', code: 'CAN' },
  { name: 'Leo', code: 'LEO' },
  { name: 'Virgo', code: 'VIR' },
  { name: 'Libra', code: 'LIB' },
  { name: 'Escorpio', code: 'ESC' },
  { name: 'Sagitario', code: 'SAG' },
  { name: 'Capricornio', code: 'CAP' },
  { name: 'Acuario', code: 'ACU' },
  { name: 'Piscis', code: 'PIS' }
])

let timeoutAux: number
const handleInputChange = (form: any) => {
  if (form.meta.valid) {
    clearTimeout(timeoutAux)
    timeoutAux = setTimeout(() => {
      emits('filters-changed', {
        cityName: GeneralUtils.clearBlankOrNull(form.values.cityName),
        countryCode: GeneralUtils.clearBlankOrNull(form.values.countryCode),
        nickName: GeneralUtils.clearBlankOrNull(form.values.nickName),
        zodiacSignCode: GeneralUtils.clearBlankOrNull(form.values.zodiacSignCode)
      })
    }, 500)
  }
}

</script>

<template>
    <div class="ModelSearchFilters">
        <div class="form">
            <Form :initialValues="initialValues" :validationSchema="schema" v-slot="form">

                <Field name="nickName" v-slot="{ field, errorMessage }" @input="handleInputChange(form)">
                    <div class="field">
                        <span class="p-float-label">
                            <InputText v-bind="field" v-model="form.values[field.name]"
                                :class="{ 'p-invalid': errorMessage }" />
                            <label for="field.name">Nick</label>
                        </span>
                        <small v-show="errorMessage" id="nick-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="countryCode" v-slot="{ field, errorMessage, handleChange }" @input="handleInputChange(form)">
                    <div class="field">
                        <span class="p-float-label">
                            <Dropdown :model-value="field.value" @update:modelValue="handleChange" :options="countries"
                                showClear optionLabel="name" optionValue="code" :class="{ 'p-invalid': errorMessage }"
                                @change="handleInputChange(form)" />
                            <label>{{ t('country') }}</label>
                        </span>
                        <small v-show="errorMessage" id="country-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="cityName" v-slot="{ field, errorMessage }" @input="handleInputChange(form)">
                    <div class="field">
                        <span class="p-float-label">
                            <InputText v-bind="field" v-model="form.values[field.name]"
                                :class="{ 'p-invalid': errorMessage }" />
                            <label>Ciudad</label>
                        </span>
                        <small v-show="errorMessage" id="city-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="zodiacSignCode" v-slot="{ field, errorMessage, handleChange }">
                    <div class="field">
                        <span class="p-float-label">
                            <Dropdown :model-value="field.value" @update:modelValue="handleChange" :options="zodiacSigns"
                                showClear optionLabel="name" optionValue="code" :class="{ 'p-invalid': errorMessage }"
                                @change="handleInputChange(form)" />
                            <label>Signo Zodiacal</label>
                        </span>
                        <small v-show="errorMessage" id="zodiac-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

            </Form>
        </div>
    </div>
</template>

<style lang="scss" scoped>
    .ModelSearchFilters {
        .form {
                padding-top: 40px;
                padding-left: 10px;
                width: 100%;
                font-size: 15px;

                .field {

                    margin-bottom: 30px;

                    .p-inputwrapper,
                    .p-inputtext {
                        width: 100%;

                        border-radius: 0;
                    }
                }

            }
    }

</style>
